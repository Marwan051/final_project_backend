package auth

import (
	"context"
	"fmt"
	"strings"

	firebase "firebase.google.com/go/v4"
	fbauth "firebase.google.com/go/v4/auth"
)

// Verifier verifies Firebase ID tokens and returns a simplified verified token
type Verifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (*VerifiedToken, error)
}

// VerifiedToken contains the UID and claims extracted from a verified token
type VerifiedToken struct {
	UID    string
	Claims map[string]any
}

// FirebaseVerifier implements Verifier using the Firebase Admin SDK
type FirebaseVerifier struct {
	client *fbauth.Client
}

// NewFirebaseVerifier initializes Firebase App using application default credentials
// and returns a FirebaseVerifier.
func NewFirebaseVerifier(ctx context.Context) (*FirebaseVerifier, error) {
	return NewFirebaseVerifierWithProjectID(ctx, "")
}

// NewFirebaseVerifierWithProjectID initializes Firebase App with ADC credentials and,
// when provided, an explicit project ID required by some local ADC user-credential flows.
func NewFirebaseVerifierWithProjectID(ctx context.Context, projectID string) (*FirebaseVerifier, error) {
	var cfg *firebase.Config
	if projectID := strings.TrimSpace(projectID); projectID != "" {
		cfg = &firebase.Config{ProjectID: projectID}
	}

	app, err := firebase.NewApp(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("firebase.NewApp: %w", err)
	}
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("app.Auth: %w (set FIREBASE_PROJECT_ID or GOOGLE_CLOUD_PROJECT when using local ADC user credentials)", err)
	}
	return &FirebaseVerifier{client: client}, nil
}

// VerifyIDToken verifies an ID token and returns the UID and claims
func (v *FirebaseVerifier) VerifyIDToken(ctx context.Context, idToken string) (*VerifiedToken, error) {
	tok, err := v.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("VerifyIDToken: %w", err)
	}
	return &VerifiedToken{UID: tok.UID, Claims: tok.Claims}, nil
}
