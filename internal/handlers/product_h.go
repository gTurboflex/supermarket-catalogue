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
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	rows, err := database.DB.Query(`
		SELECT id, name, price, stock, image, category_id, owner_id, supermarket_id, barcode, unit, unit_price, last_updated, created_at
		FROM products
		ORDER BY id
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var image, barcode, unit sql.NullString
		var unitPrice sql.NullFloat64
		var lastUpdated, createdAt sql.NullTime
		var ownerID, supermarketID sql.NullInt64

		if err := rows.Scan(
			&p.ID, &p.Name, &p.Price, &p.Stock,
			&image, &p.CategoryID, &ownerID, &supermarketID,
			&barcode, &unit, &unitPrice, &lastUpdated, &createdAt,
		); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

		products = append(products, p)
	}

	var total int
	err = database.DB.QueryRow(`SELECT COUNT(*) FROM products`).Scan(&total)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"products": products,
		"total":    total,
		"page":     page,
		"limit":    limit,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO products (name, price, stock, image, category_id, owner_id, supermarket_id, barcode, unit, unit_price)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
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
		SELECT id, name, price, stock, image, category_id, owner_id, supermarket_id, barcode, unit, unit_price, last_updated
		FROM products 
		WHERE id = $1
	`

	var p models.Product
	var image sql.NullString
	var barcode sql.NullString
	var unit sql.NullString
	var unitPrice sql.NullFloat64
	var lastUpdated sql.NullTime
	var ownerID sql.NullInt64
	var supermarketID sql.NullInt64

	err = database.DB.QueryRow(query, id).Scan(
		&p.ID, &p.Name, &p.Price, &p.Stock,
		&image, &p.CategoryID, &ownerID, &supermarketID,
		&barcode, &unit, &unitPrice, &lastUpdated,
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

	var product models.Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `
        UPDATE products 
        SET name = $1, price = $2, stock = $3, image = $4, category_id = $5, owner_id = $6, supermarket_id = $7, barcode = $8, unit = $9, unit_price = $10
        WHERE id = $11
        RETURNING id
    `

	err = database.DB.QueryRow(query,
		product.Name, product.Price, product.Stock, product.Image,
		product.CategoryID, product.OwnerID, product.SupermarketID, product.Barcode, product.Unit, product.UnitPrice, id,
	).Scan(&product.ID)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Product not found",
		})
		return
	}

	product.ID = id
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
