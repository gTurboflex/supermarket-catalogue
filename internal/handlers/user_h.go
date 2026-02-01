package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
)

// User представляет сотрудника системы (Админа или Модератора)
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"` // admin (продукты) или moderator (цены)
}

var (
	// requirement: Safe data access. Используем Mutex, чтобы при одновременном
	// обращении к списку сотрудников сервер не упал.
	userMu sync.RWMutex

	// Имитация базы данных сотрудников магазина
	staff = []User{
		{ID: 1, Name: "Aisultan", Role: "Admin"},   // Управляет списком товаров
		{ID: 2, Name: "Yergun", Role: "Moderator"}, // Управляет ценами в магазинах
	}
)

// GetUsersHandler возвращает список персонала.
// этот хендлер нужен для админ-панели,
// чтобы видеть, кто имеет доступ к редактированию каталога и цен.
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	userMu.RLock() // Блокировка на чтение
	defer userMu.RUnlock()

	w.Header().Set("Content-Type", "application/json")

	// requirement: JSON output
	if err := json.NewEncoder(w).Encode(staff); err != nil {
		http.Error(w, "Ошибка кодирования JSON", http.StatusInternalServerError)
		return
	}
}
