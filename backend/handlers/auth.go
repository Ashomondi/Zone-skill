package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"zone-skill/backend/services"
	"zone-skill/backend/utils"
)

type AuthHandler struct {
	authService *services.AuthService
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"fullName"`
}

type AuthResponse struct {
	Message string      `json:"message"`
	User    interface{} `json:"user,omitempty"`
	Token   string      `json:"token,omitempty"`
}

type GoogleAuthRequest struct {
	Credential string `json:"credential"`
	IDToken    string `json:"idToken"`
}

// NewAuthHandler creates a new authentication handler.
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles POST /api/auth/register.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteRequestError(w, r, &utils.AppError{Category: utils.CategoryValidation, UserMsg: "Method not allowed.", StatusCode: http.StatusMethodNotAllowed})
		return
	}

	var request AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.WriteRequestError(w, r, utils.NewValidationError("Invalid request body.", nil))
		return
	}

	user, token, err := h.authService.Register(
		r.Context(),
		request.Email,
		request.Password,
		request.FullName,
	)

	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidEmail):
			utils.WriteRequestError(w, r, utils.NewValidationError("Invalid email address.", nil))

		case errors.Is(err, services.ErrInvalidPassword):
			utils.WriteRequestError(w, r, utils.NewValidationError("Password must be at least 8 characters.", nil))

		case errors.Is(err, services.ErrUserAlreadyExists):
			utils.WriteRequestError(w, r, utils.NewConflictError("An account with that email already exists."))

		default:
			utils.WriteRequestError(w, r, utils.NewDatabaseError(err, map[string]string{"operation": "register_user"}))
		}

		return
	}

	utils.WriteSuccess(w, http.StatusCreated, AuthResponse{
		Message: "account created successfully",
		User:    user,
		Token:   token,
	})
}

// Login handles POST /api/auth/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteRequestError(w, r, &utils.AppError{Category: utils.CategoryValidation, UserMsg: "Method not allowed.", StatusCode: http.StatusMethodNotAllowed})
		return
	}

	var request AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.WriteRequestError(w, r, utils.NewValidationError("Invalid request body.", nil))
		return
	}

	user, token, err := h.authService.Login(
		r.Context(),
		request.Email,
		request.Password,
	)

	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			utils.WriteRequestError(w, r, utils.NewAuthError("Invalid email or password.", http.StatusUnauthorized))
			return
		}

		utils.WriteRequestError(w, r, utils.NewDatabaseError(err, map[string]string{"operation": "login_user"}))
		return
	}

	utils.WriteSuccess(w, http.StatusOK, AuthResponse{
		Message: "login successful",
		User:    user,
		Token:   token,
	})
}

// GoogleLogin handles POST /api/auth/google. The frontend sends the ID token
// returned by Google Identity Services.
func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteRequestError(w, r, &utils.AppError{Category: utils.CategoryValidation, UserMsg: "Method not allowed.", StatusCode: http.StatusMethodNotAllowed})
		return
	}

	var request GoogleAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.WriteRequestError(w, r, utils.NewValidationError("Invalid request body.", nil))
		return
	}

	credential := request.Credential
	if credential == "" {
		credential = request.IDToken
	}
	if credential == "" {
		utils.WriteRequestError(w, r, utils.NewValidationError("Google credential is required.", nil))
		return
	}

	user, token, err := h.authService.LoginWithGoogle(r.Context(), credential)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrGoogleAuthUnavailable):
			utils.WriteRequestError(w, r, utils.NewUpstreamError("Google sign-in is not configured yet.", err, map[string]string{"provider": "google"}))

		case errors.Is(err, services.ErrInvalidGoogleToken):
			utils.WriteRequestError(w, r, utils.NewAuthError("Google sign-in could not be verified. Please try again.", http.StatusUnauthorized))

		case errors.Is(err, services.ErrGoogleEmailUnverified):
			utils.WriteRequestError(w, r, utils.NewAuthError("Your Google email address must be verified before signing in.", http.StatusForbidden))

		case errors.Is(err, services.ErrGoogleAccountConflict):
			utils.WriteRequestError(w, r, utils.NewConflictError("That Google account is already linked to a different account."))

		case errors.Is(err, services.ErrAccountDisabled):
			utils.WriteRequestError(w, r, utils.NewAuthError("This account is disabled.", http.StatusForbidden))

		default:
			utils.WriteRequestError(w, r, utils.NewDatabaseError(err, map[string]string{"operation": "google_login"}))
		}
		return
	}

	utils.WriteSuccess(w, http.StatusOK, AuthResponse{
		Message: "google login successful",
		User:    user,
		Token:   token,
	})
}
