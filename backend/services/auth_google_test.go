package services

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"zone-skill/backend/repositories"
	"zone-skill/backend/utils"
)

type fakeGoogleVerifier struct {
	identity *GoogleIdentity
	err      error
}

func (f *fakeGoogleVerifier) Verify(context.Context, string) (*GoogleIdentity, error) {
	return f.identity, f.err
}

type fakeUserRepo struct {
	byEmail  map[string]*repositories.User
	byGoogle map[string]*repositories.User
	nextID   int
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{
		byEmail:  map[string]*repositories.User{},
		byGoogle: map[string]*repositories.User{},
	}
}

func (f *fakeUserRepo) next() string {
	f.nextID++
	return fmt.Sprintf("user-%d", f.nextID)
}

func (f *fakeUserRepo) Create(_ context.Context, u *repositories.User) (*repositories.User, error) {
	if _, exists := f.byEmail[u.Email]; exists {
		return nil, repositories.ErrUserEmailTaken
	}
	out := *u
	out.ID = f.next()
	out.AuthProvider = "local"
	f.byEmail[out.Email] = &out
	return &out, nil
}

func (f *fakeUserRepo) CreateGoogleUser(_ context.Context, u *repositories.User) (*repositories.User, error) {
	if _, exists := f.byEmail[u.Email]; exists {
		return nil, repositories.ErrUserEmailTaken
	}
	if _, exists := f.byGoogle[u.GoogleID]; exists {
		return nil, repositories.ErrGoogleAccountConflict
	}
	out := *u
	out.ID = f.next()
	out.AuthProvider = "google"
	out.IsActive = true
	f.byEmail[out.Email] = &out
	f.byGoogle[out.GoogleID] = &out
	return &out, nil
}

func (f *fakeUserRepo) GetByEmail(_ context.Context, email string) (*repositories.User, error) {
	if u, ok := f.byEmail[email]; ok {
		return u, nil
	}
	return nil, repositories.ErrUserNotFound
}

func (f *fakeUserRepo) GetByGoogleID(_ context.Context, googleID string) (*repositories.User, error) {
	if u, ok := f.byGoogle[googleID]; ok {
		return u, nil
	}
	return nil, repositories.ErrUserNotFound
}

func (f *fakeUserRepo) LinkGoogleAccount(_ context.Context, userID, googleID, avatarURL string) (*repositories.User, error) {
	for _, u := range f.byEmail {
		if u.ID != userID {
			continue
		}
		if u.GoogleID != "" && u.GoogleID != googleID {
			return nil, repositories.ErrGoogleAccountConflict
		}
		if other, ok := f.byGoogle[googleID]; ok && other.ID != userID {
			return nil, repositories.ErrGoogleAccountConflict
		}
		u.GoogleID = googleID
		if avatarURL != "" {
			u.AvatarURL = avatarURL
		}
		if u.AuthProvider == "local" {
			u.AuthProvider = "both"
		}
		f.byGoogle[googleID] = u
		return u, nil
	}
	return nil, repositories.ErrUserNotFound
}

func newTestAuthService(repo UserRepository, verifier GoogleTokenVerifier) *AuthService {
	return NewAuthService(repo, utils.NewJWTManager("test-secret", time.Hour), verifier)
}

func TestLoginWithGoogleCreatesUser(t *testing.T) {
	repo := newFakeUserRepo()
	verifier := &fakeGoogleVerifier{identity: &GoogleIdentity{
		Subject:       "google-subject-1",
		Email:         "new@example.com",
		EmailVerified: true,
		Name:          "New User",
		Picture:       "https://example.com/avatar.png",
	}}
	service := newTestAuthService(repo, verifier)

	user, token, err := service.LoginWithGoogle(context.Background(), "id-token")
	if err != nil {
		t.Fatalf("LoginWithGoogle returned error: %v", err)
	}
	if token == "" {
		t.Fatal("expected a token")
	}
	if user.GoogleID != "google-subject-1" {
		t.Errorf("GoogleID = %q, want google-subject-1", user.GoogleID)
	}
	if user.AuthProvider != "google" {
		t.Errorf("AuthProvider = %q, want google", user.AuthProvider)
	}
	if user.PasswordHash != "" {
		t.Errorf("PasswordHash should not be returned, got %q", user.PasswordHash)
	}
}

func TestLoginWithGoogleLinksExistingEmailAccount(t *testing.T) {
	repo := newFakeUserRepo()
	existing := &repositories.User{
		ID:           "local-user",
		Email:        "existing@example.com",
		PasswordHash: "bcrypt-hash",
		FullName:     "Existing User",
		IsActive:     true,
		AuthProvider: "local",
	}
	repo.byEmail[existing.Email] = existing

	verifier := &fakeGoogleVerifier{identity: &GoogleIdentity{
		Subject:       "google-subject-2",
		Email:         "existing@example.com",
		EmailVerified: true,
	}}
	service := newTestAuthService(repo, verifier)

	user, _, err := service.LoginWithGoogle(context.Background(), "id-token")
	if err != nil {
		t.Fatalf("LoginWithGoogle returned error: %v", err)
	}
	if user.ID != existing.ID {
		t.Errorf("ID = %q, want existing account %q", user.ID, existing.ID)
	}
	if user.GoogleID != "google-subject-2" {
		t.Errorf("GoogleID = %q, want google-subject-2", user.GoogleID)
	}
	if user.AuthProvider != "both" {
		t.Errorf("AuthProvider = %q, want both", user.AuthProvider)
	}
}

func TestLoginWithGoogleRejectsUnverifiedEmail(t *testing.T) {
	service := newTestAuthService(newFakeUserRepo(), &fakeGoogleVerifier{identity: &GoogleIdentity{
		Subject:       "google-subject-3",
		Email:         "unverified@example.com",
		EmailVerified: false,
	}})

	_, _, err := service.LoginWithGoogle(context.Background(), "id-token")
	if !errors.Is(err, ErrGoogleEmailUnverified) {
		t.Fatalf("error = %v, want ErrGoogleEmailUnverified", err)
	}
}

func TestLoginWithGoogleRejectsInvalidToken(t *testing.T) {
	service := newTestAuthService(newFakeUserRepo(), &fakeGoogleVerifier{err: errors.New("bad token")})

	_, _, err := service.LoginWithGoogle(context.Background(), "bad-token")
	if !errors.Is(err, ErrInvalidGoogleToken) {
		t.Fatalf("error = %v, want ErrInvalidGoogleToken", err)
	}
}

func TestLoginWithGoogleUnavailableWithoutVerifier(t *testing.T) {
	service := newTestAuthService(newFakeUserRepo(), nil)

	_, _, err := service.LoginWithGoogle(context.Background(), "id-token")
	if !errors.Is(err, ErrGoogleAuthUnavailable) {
		t.Fatalf("error = %v, want ErrGoogleAuthUnavailable", err)
	}
}

func TestLoginWithGoogleRejectsConflictingLink(t *testing.T) {
	repo := newFakeUserRepo()
	existing := &repositories.User{
		ID:           "linked-user",
		Email:        "linked@example.com",
		PasswordHash: "bcrypt-hash",
		FullName:     "Linked User",
		IsActive:     true,
		GoogleID:     "different-google-subject",
		AuthProvider: "both",
	}
	repo.byEmail[existing.Email] = existing
	repo.byGoogle[existing.GoogleID] = existing

	service := newTestAuthService(repo, &fakeGoogleVerifier{identity: &GoogleIdentity{
		Subject:       "google-subject-4",
		Email:         "linked@example.com",
		EmailVerified: true,
	}})

	_, _, err := service.LoginWithGoogle(context.Background(), "id-token")
	if !errors.Is(err, ErrGoogleAccountConflict) {
		t.Fatalf("error = %v, want ErrGoogleAccountConflict", err)
	}
}
