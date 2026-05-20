package authctx

import "context"

type contextKey string

const (
	userUIDKey    contextKey = "userUID"
	userClaimsKey contextKey = "userClaims"
	authTokenKey  contextKey = "authToken"
)

func SetUserUID(ctx context.Context, uid string) context.Context {
	return context.WithValue(ctx, userUIDKey, uid)
}

func SetUserClaims(ctx context.Context, claims map[string]any) context.Context {
	return context.WithValue(ctx, userClaimsKey, claims)
}

func SetAuthToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, authTokenKey, token)
}

func GetUserUID(ctx context.Context) (string, bool) {
	v := ctx.Value(userUIDKey)
	s, ok := v.(string)
	return s, ok
}

func GetUserClaims(ctx context.Context) (map[string]any, bool) {
	v := ctx.Value(userClaimsKey)
	m, ok := v.(map[string]any)
	return m, ok
}

func GetAuthToken(ctx context.Context) (string, bool) {
	v := ctx.Value(authTokenKey)
	s, ok := v.(string)
	return s, ok
}
