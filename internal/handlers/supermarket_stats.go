package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	database "supermarket-catalogue/internal/repository"
)

type SupermarketStats struct {
	SupermarketID   int     `json:"supermarket_id"`
	SupermarketName string  `json:"supermarket_name"`
	ProductCount    int     `json:"product_count"`
	AvgPrice        float64 `json:"avg_price"`
	MinPrice        float64 `json:"min_price"`
	MaxPrice        float64 `json:"max_price"`
}

// GetSupermarketStats returns stats for each supermarket
// @Summary Supermarket statistics
// @Description Retrieve product statistics for each supermarket: product count, average, min, max price
// @Tags supermarkets
// @Accept json
// @Produce json
// @Success 200 {array} SupermarketStats
// @Router /supermarkets/stats [get]
func GetSupermarketStats(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT s.id, s.name, COUNT(p.id), AVG(p.price), MIN(p.price), MAX(p.price)
		FROM supermarkets s
		LEFT JOIN products p ON p.supermarket_id = s.id
		GROUP BY s.id, s.name
		ORDER BY s.id
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		http.Error(w, "database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var stats []SupermarketStats
	for rows.Next() {
		var s SupermarketStats
		var avg sql.NullFloat64
		var min sql.NullFloat64
		var max sql.NullFloat64
		err := rows.Scan(&s.SupermarketID, &s.SupermarketName, &s.ProductCount, &avg, &min, &max)
		if err != nil {
			http.Error(w, "scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if avg.Valid {
			s.AvgPrice = avg.Float64
		}
		if min.Valid {
			s.MinPrice = min.Float64
		}
		if max.Valid {
			s.MaxPrice = max.Float64
		}
		stats = append(stats, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
