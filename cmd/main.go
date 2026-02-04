package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/airlangga-hub/clothing-ecommerce/entity"
	"github.com/airlangga-hub/clothing-ecommerce/handler"
	"github.com/airlangga-hub/clothing-ecommerce/helper"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	dsn := os.Getenv("DSN")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalln("Error connecting to MySQL:", err)
	}
	defer db.Close()

	h := handler.NewHandler(db)

	// variables
	var (
		user      entity.User
		u         entity.User
		products  []entity.Product
		product   entity.Product
		cartItems []entity.CartItem
		priceStr  string
		input     string
		buf       bytes.Buffer
		price     int
	)
	w := tabwriter.NewWriter(&buf, 0, 0, 1, ' ', tabwriter.Debug)
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("\n------- Welcome to Hacktiv8 Clothing Store -------")

MainMenu:
	fmt.Println("\nMain Menu:")
	fmt.Println("1. Register")
	fmt.Println("2. Login")
	fmt.Println("3. Exit")
	fmt.Print("Your input (1/2/3): ")

	scanner.Scan()
	input = strings.TrimSpace(scanner.Text())

	switch input {
	case "1":
		goto Register
	case "2":
		goto Login
	default:
		goto Exit
	}

Register:
	fmt.Print("\nEmail: ")
	scanner.Scan()
	user.Email = strings.TrimSpace(scanner.Text())

	fmt.Print("Password: ")
	scanner.Scan()
	user.Password = strings.TrimSpace(scanner.Text())

	if err := h.CreateUser(user.Email, user.Password); err != nil {
		slog.Error(err.Error())
		goto MainMenu
	}

	user, err = h.ReadUserByEmail(user.Email)
	if err != nil {
		slog.Error(err.Error())
		goto MainMenu
	}

	fmt.Println("\nRegister success!!!!")

	goto UserMenu

Login:
	fmt.Print("\nEmail: ")
	scanner.Scan()
	user.Email = strings.TrimSpace(scanner.Text())

	fmt.Print("Password: ")
	scanner.Scan()
	user.Password = strings.TrimSpace(scanner.Text())

	u, err = h.ReadUserByEmail(user.Email)
	if err != nil {
		slog.Error(err.Error())
		goto MainMenu
	}

	if u.Password != user.Password {
		fmt.Println("\nWrong Pasword!!!!!")
		goto Login
	}

	user = u

	fmt.Println("\nLogin success!!!!")

	if user.Role == "user" {
		goto UserMenu
	} else {
		goto AdminMenu
	}

Exit:
	fmt.Println("\nGoodbye!!!!")
	return

AdminMenu:
	fmt.Println("\nAdmin Menu:")
	fmt.Println("1. Show All Products")
	fmt.Println("2. Create Product")

	scanner.Scan()
	input = strings.TrimSpace(scanner.Text())

	switch input {
	case "1":
		goto CreateProduct
	case "2":
		goto ShowAllProducts
	default:
		goto Exit
	}

UserMenu:
	fmt.Println("\nUser Menu:")
	fmt.Println("1. Show All Products")
	fmt.Println("2. Add To Cart")
	fmt.Println("3. Show Cart")
	fmt.Println("4. Create Order")
	fmt.Println("5. Exit")
	fmt.Print("Your input (1/2/3/4/5): ")

	scanner.Scan()
	input = strings.TrimSpace(scanner.Text())

	switch input {
	case "1":
		goto ShowAllProducts
	case "2":
		goto AddToCart
	case "3":
		goto ShowCart
	case "4":
		goto CreateOrders
	default:
		goto Exit
	}

	// Admin function
CreateProduct:
	fmt.Print("\nProduct name: ")
	scanner.Scan()
	product.Name = strings.TrimSpace(scanner.Text())
	if product.Name == "" {
		fmt.Println("Name cannot be empty!")
		goto CreateProduct
	}

	fmt.Print("Description: ")
	scanner.Scan()
	product.Description = strings.TrimSpace(scanner.Text())

	fmt.Print("Price (in Rupiah, e.g., 75000): ")
	scanner.Scan()
	priceStr = strings.TrimSpace(scanner.Text())

	price, err = strconv.Atoi(priceStr)
	if err != nil || price <= 0 {
		fmt.Println("Invalid price! Must be a positive number.")
		goto CreateProduct
	}

	err = h.CreateProduct(product.Name, product.Description, price*100)
	if err != nil {
		slog.Error(err.Error())
		fmt.Println(" Failed to create product. Please try again.")
		goto CreateProduct
	}

	fmt.Println("Product created successfully!!!!")

	// Users function
ShowAllProducts:
	fmt.Println("\nShowing all products.....")
	fmt.Fprintln(w, "| Name\t Description\t Price\t")

	products, err = h.ReadAllProducts()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	for _, product := range products {
		fmt.Fprintf(w, "| %s\t %s\t Rp%.2f\t\n", product.Name, product.Description, product.Price)
	}

	if err := w.Flush(); err != nil {
		slog.Error(err.Error())
		goto UserMenu
	}

	helper.PrintStdOut(&buf)

	goto UserMenu

AddToCart:
	fmt.Print("\nProduct ID: ")
	scanner.Scan()
	product.Id, err = strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		slog.Error(err.Error())
		fmt.Println("\nInvalid product ID!!!!")
		goto AddToCart
	}

	fmt.Print("Quantity: ")
	scanner.Scan()
	product.Quantity, err = strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		slog.Error(err.Error())
		fmt.Println("\nInvalid quantity!!!!")
		goto AddToCart
	}

	err = h.CreateCartItem(user.Id, product.Id, product.Quantity)
	if err != nil {
		slog.Error(err.Error())
		goto AddToCart
	}
	
	fmt.Println("\nAdd to cart success!!!!")
	goto UserMenu

ShowCart:
	cartItems, err = h.ReadCartItemsByUserID(user.Id)
	if len(cartItems) == 0 {
		fmt.Println("\nYour cart is empty.")
		goto ShowCart
	} else {
		fmt.Println("\nCart Contents: ")
		for _, item := range cartItems {
			fmt.Printf("- Product ID: %d, Total: %d\n", item.ProductId, item.Quantity)
		}
	}

	if err != nil {
		slog.Error(err.Error())
		fmt.Println("\nFailed to load cart. Please try again.")
		goto ShowCart
	}

CreateOrders:
	err = h.CreateOrder(entity.Order{})
	if err != nil {
		slog.Error(err.Error())
		fmt.Println("Failed to place order. Please try again.")
		return
	}
}
