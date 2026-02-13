package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"supermarket-catalogue/internal/auth"
	"supermarket-catalogue/internal/models"
	"supermarket-catalogue/internal/repository"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if user.Name == "" || user.Email == "" || user.Password == "" {
		http.Error(w, `{"error": "Name, email and password are required"}`, http.StatusBadRequest)
		return
	}

	if user.Role == "" {
		user.Role = "user"
	}

	existingUser, _ := repository.GetUserByEmail(user.Email)
	if existingUser != nil {
		http.Error(w, `{"error": "User with this email already exists"}`, http.StatusConflict)
		return
	}

	if err := repository.CreateUser(&user); err != nil {
		http.Error(w, `{"error": "Failed to create user"}`, http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		http.Error(w, `{"error": "Failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	response := models.AuthResponse{
		User: &models.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, `{"error": "Invalid email or password"}`, http.StatusUnauthorized)
		return
	}

	if !auth.CheckPasswordHash(req.Password, user.Password) {
		http.Error(w, `{"error": "Invalid email or password"}`, http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		http.Error(w, `{"error": "Failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	response := models.AuthResponse{
		User: &models.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := repository.GetAllUsers()
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch users"}`, http.StatusInternalServerError)
		return
	}

	for i := range users {
		users[i].Password = ""
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		http.Error(w, `{"error": "Not authenticated"}`, http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid user ID"}`, http.StatusBadRequest)
		return
	}

	user, err := repository.GetUserByID(userID)
	if err != nil {
		http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
