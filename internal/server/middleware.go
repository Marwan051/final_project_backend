package server

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	authpkg "github.com/Marwan051/final_project_backend/internal/auth"
	authctx "github.com/Marwan051/final_project_backend/internal/authctx"
	"github.com/Marwan051/final_project_backend/internal/utils"
)

type Middleware func(http.Handler) http.Handler

// Chain applies middlewares in order (first middleware is outermost)
func ChainMiddleware(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logging logs each incoming request with method, path, status, and duration
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		log.Printf("%s %s %d %v", r.Method, r.URL.Path, rw.status, time.Since(start))
	})
}

// PanicRecover recovers from panics and returns a 500 error
func PanicRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Caught panic: %v\nStack trace: %s", err, debug.Stack())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// CORS handles Cross-Origin Resource Sharing
func Headers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // TODO: configure allowed origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Auth returns a middleware that validates Supabase access tokens using the provided Verifier.
func Auth(verifier authpkg.Verifier) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if utils.Cfg.ENV == "dev" && utils.Cfg.DisableAuth {
				// Inject a mock dev user to prevent downstream panics
				ctx := authctx.SetUserUID(r.Context(), "dev-user-mock-uid")
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing Authorization header", http.StatusUnauthorized)
				return
			}
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "invalid Authorization header", http.StatusUnauthorized)
				return
			}
			token := parts[1]

			vt, err := verifier.VerifyIDToken(r.Context(), token)
			if err != nil {
				log.Printf("token verify failed: %v", err)
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := authctx.SetUserUID(r.Context(), vt.UID)
			ctx = authctx.SetUserClaims(ctx, vt.Claims)
			ctx = authctx.SetAuthToken(ctx, token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserUID extracts the authenticated user ID from context.
func GetUserUID(ctx context.Context) (string, bool) {
	return authctx.GetUserUID(ctx)
}

// GetUserClaims extracts the verified claims from context.
func GetUserClaims(ctx context.Context) (map[string]interface{}, bool) {
	return authctx.GetUserClaims(ctx)
}

// GetAuthToken extracts the verified bearer token from context.
func GetAuthToken(ctx context.Context) (string, bool) {
	return authctx.GetAuthToken(ctx)
}

func claimBool(claims map[string]any, key string) (bool, bool) {
	value, ok := claims[key]
	if !ok {
		return false, false
	}
	result, ok := value.(bool)
	return result, ok
}

func metadataString(claims map[string]any, container, key string) (string, bool) {
	raw, ok := claims[container]
	if !ok {
		return "", false
	}
	group, ok := raw.(map[string]any)
	if !ok {
		return "", false
	}
	value, ok := group[key]
	if !ok {
		return "", false
	}
	result, ok := value.(string)
	return result, ok
}

func metadataBool(claims map[string]any, container, key string) (bool, bool) {
	raw, ok := claims[container]
	if !ok {
		return false, false
	}
	group, ok := raw.(map[string]any)
	if !ok {
		return false, false
	}
	return claimBool(group, key)
}

func hasRole(claims map[string]any, role string) bool {
	if claimRole, ok := claims["role"].(string); ok && claimRole == role {
		return true
	}
	if metadataRole, ok := metadataString(claims, "app_metadata", "role"); ok && metadataRole == role {
		return true
	}
	return false
}

// RequirePremium enforces that the user has a premium flag or premium role in Supabase metadata.
func RequirePremium(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := GetUserClaims(r.Context())
		if !ok {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if premium, ok := claims["premium"].(bool); ok && premium {
			next.ServeHTTP(w, r)
			return
		}
		if premium, ok := metadataBool(claims, "app_metadata", "premium"); ok && premium {
			next.ServeHTTP(w, r)
			return
		}
		if premiumRole, ok := metadataString(claims, "app_metadata", "role"); ok && premiumRole == "premium" {
			next.ServeHTTP(w, r)
			return
		}
		if role, ok := claims["role"].(string); ok && role == "premium" {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, "forbidden", http.StatusForbidden)
	})
}

// RequireAdmin enforces that the user has an admin role or admin flag in Supabase metadata.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := GetUserClaims(r.Context())
		if !ok {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if hasRole(claims, "admin") {
			next.ServeHTTP(w, r)
			return
		}
		if admin, ok := claimBool(claims, "admin"); ok && admin {
			next.ServeHTTP(w, r)
			return
		}
		if admin, ok := metadataBool(claims, "app_metadata", "admin"); ok && admin {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, "forbidden", http.StatusForbidden)
	})
}
