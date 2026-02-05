package entity

import "time"

type User struct {
	Id       int
	Email    string
	Password string
	Role     string
}

type Product struct {
	Id          int
	Name        string
	Description string
	Price       float32
	Quantity    int
	Stock       int
}

type CartItem struct {
	Id        int
	UserId    int
	ProductId int
	Quantity  int
}

type OrderItem struct {
	Id        int
	OrderId   int
	ProductId int
	Quantity  int
}

type Order struct {
	Id         int
	UserId     int
	TotalPrice float32
	CreatedAt  time.Time
	Products   []Product
}
