package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	database "supermarket-catalogue/internal/repository"

	"github.com/gorilla/mux"
)

type compareRow struct {
	ProductID       int      `json:"product_id"`
	Name            string   `json:"name"`
	Price           float64  `json:"price"`
	UnitPrice       *float64 `json:"unit_price,omitempty"`
	Unit            string   `json:"unit,omitempty"`
	SupermarketID   *int     `json:"supermarket_id,omitempty"`
	SupermarketName string   `json:"supermarket_name,omitempty"`
	LastUpdated     *string  `json:"last_updated,omitempty"`
}

type compareResponse struct {
	Barcode string       `json:"barcode"`
	Results []compareRow `json:"results"`
	Best    *compareRow  `json:"best,omitempty"`
}

func CompareByBarcode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code := vars["barcode"]
	if code == "" {
		http.Error(w, "barcode required", http.StatusBadRequest)
		return
	}

	query := `
		SELECT p.id, p.name, p.price, p.unit_price, p.unit, p.supermarket_id, s.name, p.last_updated
		FROM products p
		LEFT JOIN supermarkets s ON p.supermarket_id = s.id
		WHERE p.barcode = $1
		ORDER BY p.unit_price IS NULL, p.unit_price ASC, p.price ASC
	`

	rows, err := database.DB.Query(query, code)
	if err != nil {
		http.Error(w, "database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var resp compareResponse
	resp.Barcode = code
	var bestIndex = -1
	var bestUnit float64
	var haveUnit bool
	for rows.Next() {
		var id int
		var name string
		var price float64
		var unitPrice sql.NullFloat64
		var unit sql.NullString
		var supermarketID sql.NullInt64
		var supermarketName sql.NullString
		var lastUpdated sql.NullTime

		if err := rows.Scan(&id, &name, &price, &unitPrice, &unit, &supermarketID, &supermarketName, &lastUpdated); err != nil {
			http.Error(w, "scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		row := compareRow{
			ProductID: id,
			Name:      name,
			Price:     price,
		}
		if unitPrice.Valid {
			up := unitPrice.Float64
			row.UnitPrice = &up
		}
		if unit.Valid {
			row.Unit = unit.String
		}
		if supermarketID.Valid {
			sid := int(supermarketID.Int64)
			row.SupermarketID = &sid
		}
		if supermarketName.Valid {
			row.SupermarketName = supermarketName.String
		}
		if lastUpdated.Valid {
			lu := lastUpdated.Time.UTC().Format("2006-01-02T15:04:05Z")
			row.LastUpdated = &lu
		}

		resp.Results = append(resp.Results, row)

		if row.UnitPrice != nil {
			if !haveUnit || *row.UnitPrice < bestUnit {
				haveUnit = true
				bestUnit = *row.UnitPrice
				bestIndex = len(resp.Results) - 1
			}
		} else if !haveUnit {
			if bestIndex == -1 || row.Price < resp.Results[bestIndex].Price {
				bestIndex = len(resp.Results) - 1
			}
		}
	}

	if len(resp.Results) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "no offers found for barcode"})
		return
	}

	if bestIndex >= 0 {
		resp.Best = &resp.Results[bestIndex]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
