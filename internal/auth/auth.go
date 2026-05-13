package auth

import (
	"context"
	"fmt"
	"strings"

	supabase "github.com/supabase-community/supabase-go"
)

// Verifier verifies Supabase access tokens and returns a simplified verified token.
type Verifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (*VerifiedToken, error)
}

// VerifiedToken contains the UID and claims extracted from a verified token.
type VerifiedToken struct {
	UID    string
	Claims map[string]any
}

// SupabaseVerifier implements Verifier using the Supabase Go SDK.
type SupabaseVerifier struct {
	client *supabase.Client
}

// NewSupabaseVerifier creates a verifier backed by the Supabase client.
func NewSupabaseVerifier(baseURL, apiKey string) (*SupabaseVerifier, error) {
	baseURL = strings.TrimSpace(baseURL)
	apiKey = strings.TrimSpace(apiKey)
	if baseURL == "" {
		return nil, fmt.Errorf("supabase base url is required")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("supabase secret key is required")
	}
	client, err := supabase.NewClient(strings.TrimRight(baseURL, "/"), apiKey, nil)
	if err != nil {
		return nil, fmt.Errorf("supabase.NewClient: %w", err)
	}
	return &SupabaseVerifier{client: client}, nil
}

// VerifyIDToken verifies a Supabase access token via the GoTrue client.
func (v *SupabaseVerifier) VerifyIDToken(ctx context.Context, idToken string) (*VerifiedToken, error) {
	_ = ctx
	authClient := v.client.Auth.WithToken(idToken)
	userResp, err := authClient.GetUser()
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	claims := map[string]any{
		"sub":           userResp.ID.String(),
		"email":         userResp.Email,
		"role":          userResp.Role,
		"aud":           userResp.Aud,
		"app_metadata":  userResp.AppMetadata,
		"user_metadata": userResp.UserMetadata,
	}

	return &VerifiedToken{UID: userResp.ID.String(), Claims: claims}, nil
}
