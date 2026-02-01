package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

var (
	userMu sync.RWMutex
	staff  = []User{
		{ID: 1, Name: "Aisultan", Role: "Admin"},
		{ID: 2, Name: "Yergun", Role: "Moderator"},
	}
)

// GetUsersHandler returns team members
// @Summary Get team members
// @Description Retrieve list of team members working on the project
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} User
// @Router /users [get]
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	userMu.RLock()
	defer userMu.RUnlock()

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(staff); err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}
