package main

import (
	"database/sql"
)

type Handler struct {
	DB *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{DB: db}
}

// create user
func (h *Handler) CreateUser(email, password string) error {
	_, err := h.DB.Exec(
		`INSERT INTO users
			(email, password)
		VALUES
			(?, ?);`,
		email, password,
	)
	
	if err != nil {
		return err
	}
	
	return nil
}

// read user by email
func (h *Handler) ReadUserByEmail(email string) (User, error) {
	var user User
	
	if err := h.DB.QueryRow(
		`SELECT
			id,
			email,
			password,
			role
		FROM users
		WHERE email = ?;`,
		email,
	).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.Role,
	); err != nil {
		return User{}, err
	}
	
	return user, nil
}

// create product
func (h *Handler) CreateProduct (name, description string, price int) error {	
	_, err := h.DB.Exec(
		`INSERT INTO products
			(name, description, price)
		VALUES
			(?, ?, ?);`,
		name, description, price,
	)
	
	if err != nil {
		return err
	}
	
	return nil
}

// read all products


// create cart item


// read cart items by user id


// delete cart items by user id


// create order by user id


// read orders by user id


// create order items

