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
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
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
		Username: req.Username,
		Password: req.Password,
	}

	token, err := h.loginHandler.Handle(r.Context(), cmd)
	if err != nil {
		h.logger.Warn("login failed", zap.Error(err))
		h.respondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	h.respondSuccess(w, http.StatusOK, LoginResponse{Token: token})
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
		Username: req.Username,
		Password: req.Password,
	}

	if err := h.registerHandler.Handle(r.Context(), cmd); err != nil {
		h.logger.Warn("registration failed", zap.Error(err))
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.respondSuccess(w, http.StatusCreated, RegisterResponse{Message: "user registered successfully"})
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
