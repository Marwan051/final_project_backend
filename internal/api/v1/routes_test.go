package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	authctx "github.com/Marwan051/final_project_backend/internal/authctx"
	"github.com/Marwan051/final_project_backend/internal/utils"
)

func TestRequireAdminDeniesWhenAuthDisabled(t *testing.T) {
	originalEnv := utils.Cfg.ENV
	originalDisableAuth := utils.Cfg.DisableAuth
	t.Cleanup(func() {
		utils.Cfg.ENV = originalEnv
		utils.Cfg.DisableAuth = originalDisableAuth
	})

	utils.Cfg.ENV = "dev"
	utils.Cfg.DisableAuth = true

	handler := requireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/routing/reload-prefix-times", nil)
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, recorder.Code)
	}
}

func TestRequireAdminAllowsAdminClaims(t *testing.T) {
	originalEnv := utils.Cfg.ENV
	originalDisableAuth := utils.Cfg.DisableAuth
	t.Cleanup(func() {
		utils.Cfg.ENV = originalEnv
		utils.Cfg.DisableAuth = originalDisableAuth
	})

	utils.Cfg.ENV = "prod"
	utils.Cfg.DisableAuth = false

	handler := requireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/routing/reload-prefix-times", nil)
	request = request.WithContext(authctx.SetUserClaims(request.Context(), map[string]any{"role": "admin"}))
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
}
