package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"zone-skill/backend/services"
	"zone-skill/backend/utils"
)

type stubGoogleVerifier struct {
	identity *services.GoogleIdentity
	err      error
}

func (s *stubGoogleVerifier) Verify(context.Context, string) (*services.GoogleIdentity, error) {
	return s.identity, s.err
}

func googleLoginRequest(body string) *http.Request {
	return httptest.NewRequest(http.MethodPost, "/api/auth/google", strings.NewReader(body))
}

func newGoogleTestHandler(verifier services.GoogleTokenVerifier) *AuthHandler {
	return NewAuthHandler(services.NewAuthService(nil, utils.NewJWTManager("test-secret", time.Hour), verifier))
}

func TestGoogleLoginRequiresCredential(t *testing.T) {
	handler := newGoogleTestHandler(nil)
	rec := httptest.NewRecorder()

	handler.GoogleLogin(rec, googleLoginRequest(`{}`))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestGoogleLoginUnavailableWhenNotConfigured(t *testing.T) {
	handler := newGoogleTestHandler(nil)
	rec := httptest.NewRecorder()

	handler.GoogleLogin(rec, googleLoginRequest(`{"credential":"token"}`))

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
}

func TestGoogleLoginRejectsInvalidToken(t *testing.T) {
	handler := newGoogleTestHandler(&stubGoogleVerifier{err: errors.New("invalid")})
	rec := httptest.NewRecorder()

	handler.GoogleLogin(rec, googleLoginRequest(`{"credential":"token"}`))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestGoogleLoginRejectsUnverifiedEmail(t *testing.T) {
	handler := newGoogleTestHandler(&stubGoogleVerifier{identity: &services.GoogleIdentity{
		Subject:       "google-subject",
		Email:         "user@example.com",
		EmailVerified: false,
	}})
	rec := httptest.NewRecorder()

	handler.GoogleLogin(rec, googleLoginRequest(`{"credential":"token"}`))

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}
