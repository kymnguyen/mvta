package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/kymnguyen/mvta/apps/backend/auth-svc/internal/application/command"
	"github.com/kymnguyen/mvta/apps/backend/auth-svc/internal/application/service"
)

type AuthHandler struct {
	loginHandler    *service.LoginHandler
	registerHandler *service.RegisterUserHandler
	logger          *zap.Logger
}

func NewAuthHandler(
	loginHandler *service.LoginHandler,
	registerHandler *service.RegisterUserHandler,
	logger *zap.Logger,
) *AuthHandler {
	return &AuthHandler{
		loginHandler:    loginHandler,
		registerHandler: registerHandler,
		logger:          logger,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

type LoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

type VerifyResponse struct {
	User UserInfo `json:"user"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	cmd := &command.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	}

	token, user, err := h.loginHandler.Handle(r.Context(), cmd)
	if err != nil {
		h.logger.Warn("login failed", zap.Error(err))
		h.respondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	h.respondSuccess(w, http.StatusOK, LoginResponse{
		Token: token,
		User: UserInfo{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  user.GetRole(),
		},
	})
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	cmd := &command.RegisterUserCommand{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	if err := h.registerHandler.Handle(r.Context(), cmd); err != nil {
		h.logger.Warn("registration failed", zap.Error(err))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusCreated, RegisterResponse{Message: "user registered successfully"})
}

func (h *AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// TODO: Extract user from JWT token in Authorization header
	// For now, return a placeholder
	h.respondSuccess(w, http.StatusOK, VerifyResponse{
		User: UserInfo{
			ID:    "user-id",
			Email: "user@example.com",
			Name:  "User",
			Role:  "operator",
		},
	})
}

func (h *AuthHandler) respondSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *AuthHandler) respondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Message: message})
}
