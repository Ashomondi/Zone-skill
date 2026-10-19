package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"zone-skill/backend/repositories"
	"zone-skill/backend/utils"
)

var (
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrInvalidEmail          = errors.New("invalid email")
	ErrInvalidPassword       = errors.New("password must be at least 8 characters")
	ErrAccountDisabled       = errors.New("account is disabled")
	ErrGoogleAuthUnavailable = errors.New("google authentication is not configured")
	ErrInvalidGoogleToken    = errors.New("invalid google token")
	ErrGoogleEmailUnverified = errors.New("google email is not verified")
	ErrGoogleAccountConflict = errors.New("google account is already linked to another user")
)

// GoogleIdentity is the verified subset of a Google ID token's claims that
// this application relies on.
type GoogleIdentity struct {
	Subject       string
	Email         string
	EmailVerified bool
	Name          string
	Picture       string
}

// GoogleTokenVerifier validates a Google ID token and returns its identity
// claims. It exists as an interface so handlers/services can be tested
// without contacting Google.
type GoogleTokenVerifier interface {
	Verify(ctx context.Context, idToken string) (*GoogleIdentity, error)
}

// UserRepository defines the database operations required by authentication.
type UserRepository interface {
	Create(ctx context.Context, user *repositories.User) (*repositories.User, error)
	CreateGoogleUser(ctx context.Context, user *repositories.User) (*repositories.User, error)
	GetByEmail(ctx context.Context, email string) (*repositories.User, error)
	GetByGoogleID(ctx context.Context, googleID string) (*repositories.User, error)
	LinkGoogleAccount(ctx context.Context, userID, googleID, avatarURL string) (*repositories.User, error)
}

// AuthService handles authentication business logic.
type AuthService struct {
	userRepo       UserRepository
	jwtManager     *utils.JWTManager
	googleVerifier GoogleTokenVerifier
}

// NewAuthService creates a new authentication service. The optional Google
// verifier enables Google sign-in; when omitted, the Google endpoint reports
// that the provider is not configured.
func NewAuthService(
	userRepo UserRepository,
	jwtManager *utils.JWTManager,
	googleVerifier ...GoogleTokenVerifier,
) *AuthService {
	var verifier GoogleTokenVerifier
	if len(googleVerifier) > 0 {
		verifier = googleVerifier[0]
	}
	return &AuthService{
		userRepo:       userRepo,
		jwtManager:     jwtManager,
		googleVerifier: verifier,
	}
}

// Register creates a new user account and returns the user plus an
// authenticated session token.
func (s *AuthService) Register(
	ctx context.Context,
	email string,
	password string,
	fullName string,
) (*repositories.User, string, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	if !strings.Contains(email, "@") {
		return nil, "", ErrInvalidEmail
	}

	if len(password) < 8 {
		return nil, "", ErrInvalidPassword
	}

	// Check whether the email is already registered.
	existingUser, err := s.userRepo.GetByEmail(ctx, email)

	if err == nil && existingUser != nil {
		return nil, "", ErrUserAlreadyExists
	}

	// Hash the password before storing it.
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	now := time.Now()

	created, err := s.userRepo.Create(ctx, &repositories.User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: hashedPassword,
		FullName:     strings.TrimSpace(fullName),
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
	if err != nil {
		return nil, "", err
	}

	token, err := s.jwtManager.GenerateToken(created.ID, created.Email)
	if err != nil {
		return nil, "", err
	}

	// Never return the password hash.
	created.PasswordHash = ""

	return created, token, nil
}

// Login authenticates an existing user and returns the user plus a JWT.
func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (*repositories.User, string, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	if user == nil {
		return nil, "", ErrInvalidCredentials
	}

	// Compare the supplied password against the stored bcrypt hash. Google-only
	// accounts have no password hash and must use Google sign-in.
	if user.PasswordHash == "" || !utils.CheckPassword(password, user.PasswordHash) {
		return nil, "", ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, "", ErrAccountDisabled
	}

	token, err := s.jwtManager.GenerateToken(
		user.ID,
		user.Email,
	)
	if err != nil {
		return nil, "", err
	}

	// Never return the password hash.
	user.PasswordHash = ""

	return user, token, nil
}

// LoginWithGoogle verifies a Google ID token, finds/creates/links the local
// user, and returns an application JWT.
func (s *AuthService) LoginWithGoogle(ctx context.Context, idToken string) (*repositories.User, string, error) {
	if s.googleVerifier == nil {
		return nil, "", ErrGoogleAuthUnavailable
	}

	identity, err := s.googleVerifier.Verify(ctx, idToken)
	if err != nil || identity == nil || strings.TrimSpace(identity.Subject) == "" {
		return nil, "", ErrInvalidGoogleToken
	}

	email := strings.ToLower(strings.TrimSpace(identity.Email))
	if email == "" || !identity.EmailVerified {
		return nil, "", ErrGoogleEmailUnverified
	}

	user, err := s.userRepo.GetByGoogleID(ctx, identity.Subject)
	if err != nil && !errors.Is(err, repositories.ErrUserNotFound) {
		return nil, "", err
	}

	if errors.Is(err, repositories.ErrUserNotFound) {
		// No account is linked to this Google identity yet. Link it to an
		// existing account with the same verified email, or create a new one.
		existing, lookupErr := s.userRepo.GetByEmail(ctx, email)
		switch {
		case lookupErr == nil && existing != nil:
			if existing.GoogleID != "" && existing.GoogleID != identity.Subject {
				return nil, "", ErrGoogleAccountConflict
			}
			user, err = s.userRepo.LinkGoogleAccount(ctx, existing.ID, identity.Subject, identity.Picture)
			if err != nil {
				return nil, "", err
			}

		case errors.Is(lookupErr, repositories.ErrUserNotFound):
			fullName := strings.TrimSpace(identity.Name)
			if fullName == "" {
				fullName = email
			}
			user, err = s.userRepo.CreateGoogleUser(ctx, &repositories.User{
				Email:        email,
				FullName:     fullName,
				GoogleID:     identity.Subject,
				AuthProvider: "google",
				AvatarURL:    identity.Picture,
				IsActive:     true,
			})
			if err != nil {
				// Another request may have created the account between our
				// lookups and this insert. Recover by re-reading or linking.
				switch {
				case errors.Is(err, repositories.ErrUserEmailTaken):
					existing, lookupErr := s.userRepo.GetByEmail(ctx, email)
					if lookupErr != nil {
						return nil, "", err
					}
					if existing.GoogleID != "" && existing.GoogleID != identity.Subject {
						return nil, "", ErrGoogleAccountConflict
					}
					user, err = s.userRepo.LinkGoogleAccount(ctx, existing.ID, identity.Subject, identity.Picture)
					if err != nil {
						return nil, "", err
					}
				case errors.Is(err, repositories.ErrGoogleAccountConflict):
					existing, lookupErr := s.userRepo.GetByGoogleID(ctx, identity.Subject)
					if lookupErr != nil {
						return nil, "", ErrGoogleAccountConflict
					}
					user = existing
				default:
					return nil, "", err
				}
			}

		default:
			return nil, "", lookupErr
		}
	}

	if user == nil || !user.IsActive {
		return nil, "", ErrAccountDisabled
	}

	token, err := s.jwtManager.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", err
	}

	// Never return the password hash.
	user.PasswordHash = ""

	return user, token, nil
}
