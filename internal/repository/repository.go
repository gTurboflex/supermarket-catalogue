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
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role VARCHAR(20) DEFAULT 'user',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	_, err = DB.Exec(`
	INSERT INTO users (name, email, password, role) 
	VALUES ('Admin', 'admin@example.com', '$2a$10$2cJhQ6Q5bVK9z7Q8q5Z5/.JhQ6Q5bVK9z7Q8q5Z5/.JhQ6Q5bVK9z7Q8q5Z5/', 'admin')
	ON CONFLICT (email) DO NOTHING`)
	if err != nil {
		log.Println("Note: Could not insert default admin user:", err)
	}

	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS supermarkets (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		address TEXT,
		owner_id INTEGER REFERENCES users(id),
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
		owner_id INTEGER REFERENCES users(id),
		supermarket_id INTEGER REFERENCES supermarkets(id),
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
