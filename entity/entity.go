package entity

import "time"

type User struct {
	Id       int
	Email    string
	Password string
	Role     string
}

type UserReport struct {
	Id            int
	Email         string
	TotalSpending float64
}

type Product struct {
	Id          int
	Name        string
	Description string
	Price       float64
	Quantity    int
	Stock       int
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
	TotalPrice float64
	CreatedAt  time.Time
	Products   []Product
}

type InsertOrder struct {
	Id         int
	UserId     int
	TotalPrice int
	CreatedAt  time.Time
	Products   []Product
}

type StockReport struct {
	ProductId   int
	ProductName string
	Stock       int
	Label       string
}
