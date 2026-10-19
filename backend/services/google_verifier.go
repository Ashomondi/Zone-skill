package services

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/api/idtoken"
)

// GoogleIDTokenVerifier validates Google Identity Services ID tokens against
// the configured OAuth client ID using Google's published signing keys.
type GoogleIDTokenVerifier struct {
	clientID string
}

// NewGoogleIDTokenVerifier constructs a verifier for the given OAuth Web
// Client ID. The client ID is public and is the expected token audience.
func NewGoogleIDTokenVerifier(clientID string) *GoogleIDTokenVerifier {
	return &GoogleIDTokenVerifier{clientID: strings.TrimSpace(clientID)}
}

// Verify implements GoogleTokenVerifier.
func (v *GoogleIDTokenVerifier) Verify(ctx context.Context, rawToken string) (*GoogleIdentity, error) {
	if v == nil || v.clientID == "" {
		return nil, ErrGoogleAuthUnavailable
	}
	if strings.TrimSpace(rawToken) == "" {
		return nil, ErrInvalidGoogleToken
	}

	payload, err := idtoken.Validate(ctx, rawToken, v.clientID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidGoogleToken, err)
	}

	claims := payload.Claims
	subject := payload.Subject
	if subject == "" {
		subject = claimString(claims, "sub")
	}

	return &GoogleIdentity{
		Subject:       subject,
		Email:         claimString(claims, "email"),
		EmailVerified: claimBool(claims, "email_verified"),
		Name:          claimString(claims, "name"),
		Picture:       claimString(claims, "picture"),
	}, nil
}

func claimString(claims map[string]interface{}, key string) string {
	value, ok := claims[key]
	if !ok || value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return typed
	case fmt.Stringer:
		return typed.String()
	default:
		return fmt.Sprint(typed)
	}
}

func claimBool(claims map[string]interface{}, key string) bool {
	value, ok := claims[key]
	if !ok || value == nil {
		return false
	}
	switch typed := value.(type) {
	case bool:
		return typed
	case string:
		return strings.EqualFold(typed, "true")
	default:
		return false
	}
}
