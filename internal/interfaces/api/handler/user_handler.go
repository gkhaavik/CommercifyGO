package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/infrastructure/auth"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userUseCase *usecase.UserUseCase
	jwtService  *auth.JWTService
	logger      logger.Logger
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userUseCase *usecase.UserUseCase, jwtService *auth.JWTService, logger logger.Logger) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		jwtService:  jwtService,
		logger:      logger,
	}
}

// Register handles user registration
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input usecase.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.Register(input)
	if err != nil {
		h.logger.Error("Failed to register user: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate JWT token
	token, err := h.jwtService.GenerateToken(user)
	if err != nil {
		h.logger.Error("Failed to generate token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return user and token
	response := map[string]interface{}{
		"user":  user,
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login handles user login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input usecase.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.Login(input)
	if err != nil {
		h.logger.Error("Login failed: %v", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := h.jwtService.GenerateToken(user)
	if err != nil {
		h.logger.Error("Failed to generate token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return user and token
	response := map[string]interface{}{
		"user":  user,
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetProfile handles getting the user's profile
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userUseCase.GetUserByID(userID)
	if err != nil {
		h.logger.Error("Failed to get user profile: %v", err)
		http.Error(w, "Failed to get user profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateProfile handles updating the user's profile
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input usecase.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.UpdateUser(userID, input)
	if err != nil {
		h.logger.Error("Failed to update user profile: %v", err)
		http.Error(w, "Failed to update user profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// ChangePassword handles changing the user's password
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input usecase.ChangePasswordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.userUseCase.ChangePassword(userID, input); err != nil {
		h.logger.Error("Failed to change password: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password changed successfully"})
}

// ListUsers handles listing all users (admin only)
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10 // Default limit
	}

	users, err := h.userUseCase.ListUsers(offset, limit)
	if err != nil {
		h.logger.Error("Failed to list users: %v", err)
		http.Error(w, "Failed to list users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
