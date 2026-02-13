package handlers

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"supermarket-catalogue/internal/models"
	database "supermarket-catalogue/internal/repository"
	"time"

	"github.com/gorilla/mux"
)

func GetSupermarkets(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
		SELECT id, name, address, owner_id, created_at
		FROM supermarkets
		ORDER BY id
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.Supermarket
	for rows.Next() {
		var s models.Supermarket
		var addr sql.NullString
		var ownerID sql.NullInt64
		var createdAt sql.NullTime

		if err := rows.Scan(&s.ID, &s.Name, &addr, &ownerID, &createdAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.Address = addr.String
		if ownerID.Valid {
			s.OwnerID = int(ownerID.Int64)
		}
		if createdAt.Valid {
			s.CreatedAt = createdAt.Time
		}

		items = append(items, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func GetSupermarketByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	query := `
		SELECT id, name, address, owner_id, created_at
		FROM supermarkets
		WHERE id = $1
	`
	var s models.Supermarket
	var addr sql.NullString
	var ownerID sql.NullInt64
	var createdAt sql.NullTime

	err = database.DB.QueryRow(query, id).Scan(&s.ID, &s.Name, &addr, &ownerID, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.Address = addr.String
	if ownerID.Valid {
		s.OwnerID = int(ownerID.Int64)
	}
	if createdAt.Valid {
		s.CreatedAt = createdAt.Time
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func CreateSupermarket(w http.ResponseWriter, r *http.Request) {
	var s models.Supermarket
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO supermarkets (name, address, owner_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	var createdAt time.Time
	err := database.DB.QueryRow(query, s.Name, s.Address, s.OwnerID).Scan(&s.ID, &createdAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.CreatedAt = createdAt

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func UpdateSupermarket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var s models.Supermarket
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}

	query := `
		UPDATE supermarkets
		SET name = $1, address = $2, owner_id = $3
		WHERE id = $4
		RETURNING created_at
	`
	var createdAt time.Time
	err = database.DB.QueryRow(query, s.Name, s.Address, s.OwnerID, id).Scan(&createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.ID = id
	s.CreatedAt = createdAt

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func DeleteSupermarket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec(`DELETE FROM supermarkets WHERE id = $1`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func AdminPage(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("ui", "html", "admin_supermarkets.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "template parse error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title": "Admin — Supermarkets",
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = tmpl.Execute(w, data)
}
