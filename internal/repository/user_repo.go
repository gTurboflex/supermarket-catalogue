package repository

import (
	"supermarket-catalogue/internal/auth"
	"supermarket-catalogue/internal/models"
)

func CreateUser(user *models.User) error {
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (name, email, password, role) 
	          VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	err = DB.QueryRow(query, user.Name, user.Email, hashedPassword, user.Role).
		Scan(&user.ID, &user.CreatedAt)
	return err
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, name, email, password, role, created_at 
	          FROM users WHERE email = $1`
	err := DB.QueryRow(query, email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(id int) (*models.User, error) {
	var user models.User
	query := `SELECT id, name, email, role, created_at 
	          FROM users WHERE id = $1`
	err := DB.QueryRow(query, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetAllUsers() ([]models.User, error) {
	rows, err := DB.Query(`SELECT id, name, email, role, created_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
