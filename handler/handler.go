package handler

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/airlangga-hub/clothing-ecommerce/entity"
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
func (h *Handler) ReadUserByEmail(email string) (entity.User, error) {
	var user entity.User

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
		return entity.User{}, err
	}

	return user, nil
}

// create product
func (h *Handler) CreateProduct(name, description string, price, stock int) error {
	_, err := h.DB.Exec(
		`INSERT INTO products
			(name, description, price, stock)
		VALUES
			(?, ?, ?, ?);`,
		name, description, price, stock,
	)

	if err != nil {
		return err
	}

	return nil
}

// read products by product ids
func (h *Handler) ReadProductsByProductIDs(productIDs []int) ([]entity.Product, error) {
	questionMarks := make([]string, len(productIDs))
	IDs := make([]any, len(productIDs))

	for i, id := range productIDs {
		questionMarks[i] = "?"
		IDs[i] = id
	}

	query := fmt.Sprintf(
		`SELECT
			id,
			name,
			description,
			price
		FROM products
		WHERE id IN (%s)`,
		strings.Join(questionMarks, ","),
	)

	rows, err := h.DB.Query(query, IDs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]entity.Product, 0, 10)

	for rows.Next() {
		var product entity.Product
		var price int

		if err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&price,
		); err != nil {
			return nil, err
		}

		product.Price = float32(price) / 100

		products = append(products, product)
	}

	return products, nil
}

// read all products
func (h *Handler) ReadAllProducts() ([]entity.Product, error) {
	rows, err := h.DB.Query(
		`SElECT
			id,
			name,
			description,
			price,
			stock
		FROM products;`,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]entity.Product, 0, 10)

	for rows.Next() {
		var product entity.Product
		var price int

		if err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&price,
			&product.Stock,
		); err != nil {
			return nil, err
		}

		product.Price = float32(price) / 100

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// create cart item
func (h *Handler) CreateCartItem(userID, productID, quantity int) error {
	_, err := h.DB.Exec(
		`INSERT INTO cart_items
			(user_id, product_id, quantity)
		VALUES
			(?, ?, ?) AS new
		ON DUPLICATE KEY
		UPDATE quantity = cart_items.quantity + new.quantity;`,
		userID, productID, quantity,
	)

	if err != nil {
		return err
	}

	return nil
}

// read cart items by user id
func (h *Handler) ReadCartItemsByUserID(userID int) ([]entity.Product, error) {
	rows, err := h.DB.Query(
		`SELECT
			p.id,
			p.name,
			p.description,
			p.price,
			ci.quantity
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE user_id = ?;`,
		userID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]entity.Product, 0, 10)

	for rows.Next() {
		var product entity.Product
		var price int

		if err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&price,
			&product.Quantity,
		); err != nil {
			return nil, err
		}

		product.Price = float32(price) / 100

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// delete cart items by user id
func (h *Handler) DeleteCartItemsByUserID(userID int) error {
	_, err := h.DB.Exec(
		`DELETE FROM cart_items
		WHERE user_id = ?;`,
		userID,
	)

	if err != nil {
		return err
	}

	return nil
}

// create order
func (h *Handler) CreateOrder(order entity.Order) error {
	tx, err := h.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec(
		`INSERT INTO orders
			(user_id, total_price)
		VALUES
			(?, ?);`,
		order.UserId, order.TotalPrice,
	)
	if err != nil {
		return err
	}

	orderId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	order.Id = int(orderId)

	stmt, err := tx.Prepare(
		`CALL place_order_items(?, ?, ?)`,
	)
	if err != nil {
		return err
	}

	for _, product := range order.Products {
		_, err := stmt.Exec(order.Id, product.Id, product.Quantity)
		if err != nil {
			return err
		}
	}

	// delete cart items
	if err := h.DeleteCartItemsByUserID(order.UserId); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// read orders by user id
func (h *Handler) ReadOrdersByUserID(userID int) ([]entity.Order, error) {
	rows, err := h.DB.Query(
		`SELECT
			o.id,
			o.user_id,
			o.total_price,
			o.created_at,
			oi.product_id,
			oi.quantity,
		FROM orders o
		JOIN order_items oi ON o.id = oi.order_id
		WHERE o.user_id = ?;`,
		userID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mapOrders := make(map[int]entity.Order)
	productIDset := make(map[int]struct{})

	for rows.Next() {
		var order entity.Order
		var product entity.Product

		if err := rows.Scan(
			&order.Id,
			&order.UserId,
			&order.TotalPrice,
			&order.CreatedAt,
			&product.Id,
			&product.Quantity,
		); err != nil {
			return nil, err
		}

		productIDset[product.Id] = struct{}{}

		if order, exist := mapOrders[order.Id]; exist {
			order.Products = append(
				order.Products,
				entity.Product{
					Id:       product.Id,
					Quantity: product.Quantity,
				},
			)
		} else {
			mapOrders[order.Id] = entity.Order{
				Id:         order.Id,
				UserId:     order.UserId,
				TotalPrice: order.TotalPrice,
				CreatedAt:  order.CreatedAt,
				Products: []entity.Product{
					{
						Id:       product.Id,
						Quantity: product.Quantity,
					},
				},
			}
		}
	}

	productIDs := make([]int, 0, len(productIDset))
	for id := range productIDset {
		productIDs = append(productIDs, id)
	}

	products, err := h.ReadProductsByProductIDs(productIDs)
	if err != nil {
		return nil, err
	}

	mapProducts := make(map[int]entity.Product)
	for _, product := range products {
		mapProducts[product.Id] = product
	}

	orders := make([]entity.Order, 0, len(mapOrders))
	for _, order := range mapOrders {
		for i, product := range order.Products {
			if p, exist := mapProducts[product.Id]; exist {
				order.Products[i].Name = p.Name
				order.Products[i].Description = p.Description
				order.Products[i].Price = p.Price
			}
		}
		orders = append(orders, order)
	}

	return orders, nil
}
