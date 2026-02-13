package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"supermarket-catalogue/internal/repository"
)

type BasketItem struct {
	Barcode  string `json:"barcode"`
	Quantity int    `json:"quantity"`
}

type BasketRequest struct {
	Items []BasketItem `json:"items"`
}

type SupermarketTotal struct {
	SupermarketID   int      `json:"supermarket_id"`
	SupermarketName string   `json:"supermarket_name,omitempty"`
	Total           float64  `json:"total"`
	Missing         []string `json:"missing"`
	MatchedItems    int      `json:"matched_items"`
}

type BasketResponse struct {
	Results []SupermarketTotal `json:"results"`
}

func CompareBasket(w http.ResponseWriter, r *http.Request) {
	var req BasketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if len(req.Items) == 0 {
		http.Error(w, "no items provided", http.StatusBadRequest)
		return
	}

	type sm struct {
		ID   int
		Name string
	}
	supermarkets := []sm{}
	rows, err := repository.DB.Query("SELECT id, name FROM supermarkets")
	if err != nil {
		http.Error(w, "database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var s sm
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			rows.Close()
			http.Error(w, "scan error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		supermarkets = append(supermarkets, s)
	}
	rows.Close()

	if len(supermarkets) == 0 {
		http.Error(w, "no supermarkets available", http.StatusInternalServerError)
		return
	}

	type smState struct {
		Name       string
		Total      float64
		MissingMap map[string]bool
		Matched    int
	}
	state := map[int]*smState{}
	for _, s := range supermarkets {
		mm := make(map[string]bool)
		for _, it := range req.Items {
			mm[it.Barcode] = true
		}
		state[s.ID] = &smState{
			Name:       s.Name,
			Total:      0,
			MissingMap: mm,
			Matched:    0,
		}
	}

	for _, it := range req.Items {
		minPricePerSM := map[int]float64{}
		rows, err := repository.DB.Query("SELECT supermarket_id, COALESCE(unit_price, price) FROM products WHERE barcode = $1", it.Barcode)
		if err != nil {
			http.Error(w, "database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		for rows.Next() {
			var sid sql.NullInt64
			var eff sql.NullFloat64
			if err := rows.Scan(&sid, &eff); err != nil {
				rows.Close()
				http.Error(w, "scan error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if !sid.Valid || !eff.Valid {
				continue
			}
			sidInt := int(sid.Int64)
			price := eff.Float64
			if cur, ok := minPricePerSM[sidInt]; !ok || price < cur {
				minPricePerSM[sidInt] = price
			}
		}
		rows.Close()

		for _, s := range supermarkets {
			if price, ok := minPricePerSM[s.ID]; ok {
				state[s.ID].Total += price * float64(it.Quantity)
				if state[s.ID].MissingMap[it.Barcode] {
					state[s.ID].Matched++
					state[s.ID].MissingMap[it.Barcode] = false
				}
			}
		}
	}

	resp := BasketResponse{}
	for _, s := range supermarkets {
		st := state[s.ID]
		missing := []string{}
		for bc, miss := range st.MissingMap {
			if miss {
				missing = append(missing, bc)
			}
		}
		resp.Results = append(resp.Results, SupermarketTotal{
			SupermarketID:   s.ID,
			SupermarketName: st.Name,
			Total:           st.Total,
			Missing:         missing,
			MatchedItems:    st.Matched,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
