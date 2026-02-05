package repository

import (
	"database/sql"
	"fmt"
	"log"
	"supermarket-catalogue/internal/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() error {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	log.Println("✅ Database connected successfully")
	createTables()
	return nil
}

func createTables() {
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS supermarkets (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		address TEXT,
		owner_id INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal("Failed to create supermarkets table:", err)
	}

	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		price DECIMAL(10,2) NOT NULL,
		stock INTEGER NOT NULL,
		image TEXT,
		category_id INTEGER,
		owner_id INTEGER,
		supermarket_id INTEGER,
		barcode VARCHAR(100),
		unit VARCHAR(50),
		unit_price DECIMAL(10,2),
		last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal("Failed to create products table:", err)
	}

	log.Println("✅ Tables created/verified")
}
