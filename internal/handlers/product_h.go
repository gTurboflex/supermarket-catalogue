package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"supermarket-catalogue/internal/models"
	database "supermarket-catalogue/internal/repository"

	"github.com/gorilla/mux"
)

func GetProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
		SELECT id, name, price, stock, image, category_id, owner_id, supermarket_id, barcode, unit, unit_price, last_updated, created_at
		FROM products 
		ORDER BY id
	`)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var image sql.NullString
		var barcode sql.NullString
		var unit sql.NullString
		var unitPrice sql.NullFloat64
		var lastUpdated sql.NullTime
		var createdAt sql.NullTime
		var ownerID sql.NullInt64
		var supermarketID sql.NullInt64

		err := rows.Scan(
			&p.ID, &p.Name, &p.Price, &p.Stock,
			&image, &p.CategoryID, &ownerID, &supermarketID,
			&barcode, &unit, &unitPrice, &lastUpdated, &createdAt,
		)
		if err != nil {
			http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		p.Image = image.String
		p.Barcode = barcode.String
		p.Unit = unit.String
		if unitPrice.Valid {
			p.UnitPrice = unitPrice.Float64
		}
		if lastUpdated.Valid {
			p.LastUpdated = lastUpdated.Time
		}
		if createdAt.Valid {
		}
		if ownerID.Valid {
			p.OwnerID = int(ownerID.Int64)
		}
		if supermarketID.Valid {
			p.SupermarketID = int(supermarketID.Int64)
		}

		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		http.Error(w, `{"error": "Authentication required"}`, http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid user ID"}`, http.StatusBadRequest)
		return
	}

	product.OwnerID = userID

	query := `
    INSERT INTO products (name, price, stock, image, category_id, owner_id, supermarket_id, barcode, unit, unit_price, last_updated, created_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    RETURNING id
	`

	err = database.DB.QueryRow(
		query,
		product.Name,
		product.Price,
		product.Stock,
		product.Image,
		product.CategoryID,
		product.OwnerID,
		product.SupermarketID,
		product.Barcode,
		product.Unit,
		product.UnitPrice,
	).Scan(&product.ID)

	if err != nil {
		http.Error(w, "Failed to create product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func GetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT id, name, price, stock, image, category_id, owner_id, supermarket_id, barcode, unit, unit_price, last_updated, created_at
		FROM products 
		WHERE id = $1
	`

	var p models.Product
	var image sql.NullString
	var barcode sql.NullString
	var unit sql.NullString
	var unitPrice sql.NullFloat64
	var lastUpdated sql.NullTime
	var createdAt sql.NullTime
	var ownerID sql.NullInt64
	var supermarketID sql.NullInt64

	err = database.DB.QueryRow(query, id).Scan(
		&p.ID, &p.Name, &p.Price, &p.Stock,
		&image, &p.CategoryID, &ownerID, &supermarketID,
		&barcode, &unit, &unitPrice, &lastUpdated, &createdAt,
	)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Product not found",
		})
		return
	}

	p.Image = image.String
	p.Barcode = barcode.String
	p.Unit = unit.String
	if unitPrice.Valid {
		p.UnitPrice = unitPrice.Float64
	}
	if lastUpdated.Valid {
		p.LastUpdated = lastUpdated.Time
	}
	if ownerID.Valid {
		p.OwnerID = int(ownerID.Int64)
	}
	if supermarketID.Valid {
		p.SupermarketID = int(supermarketID.Int64)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"message": "Supermarket Catalogue API is running",
	})
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	userIDStr := r.Header.Get("X-User-ID")
	role := r.Header.Get("X-User-Role")
	if userIDStr == "" {
		http.Error(w, `{"error": "Authentication required"}`, http.StatusUnauthorized)
		return
	}
	userID, _ := strconv.Atoi(userIDStr)

	var ownerID int
	err = database.DB.QueryRow("SELECT owner_id FROM products WHERE id = $1", id).Scan(&ownerID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Product not found"})
		return
	}

	if role != "admin" && userID != ownerID {
		http.Error(w, `{"error": "You can only edit your own products"}`, http.StatusForbidden)
		return
	}

	var product models.Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `
        UPDATE products 
        SET name = $1, price = $2, stock = $3, image = $4, category_id = $5, supermarket_id = $6, barcode = $7, unit = $8, unit_price = $9, last_updated = CURRENT_TIMESTAMP
        WHERE id = $10
        RETURNING id
    `

	err = database.DB.QueryRow(query,
		product.Name, product.Price, product.Stock, product.Image,
		product.CategoryID, product.SupermarketID, product.Barcode, product.Unit,
		product.UnitPrice, id,
	).Scan(&product.ID)

	if err != nil {
		http.Error(w, "Failed to update product: "+err.Error(), http.StatusInternalServerError)
		return
	}

	product.ID = id
	product.OwnerID = ownerID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	userIDStr := r.Header.Get("X-User-ID")
	role := r.Header.Get("X-User-Role")
	if userIDStr == "" {
		http.Error(w, `{"error": "Authentication required"}`, http.StatusUnauthorized)
		return
	}
	userID, _ := strconv.Atoi(userIDStr)

	var ownerID int
	err = database.DB.QueryRow("SELECT owner_id FROM products WHERE id = $1", id).Scan(&ownerID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Product not found"})
		return
	}

	if role != "admin" && userID != ownerID {
		http.Error(w, `{"error": "You can only delete your own products"}`, http.StatusForbidden)
		return
	}

	query := `DELETE FROM products WHERE id = $1`
	result, err := database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Product not found",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
