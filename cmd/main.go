package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"text/tabwriter"

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

	handler := NewHandler(db)

	// variables
	var user User
	var u User
	var input string
	var buf bytes.Buffer
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

	if err := handler.CreateUser(user.Email, user.Password); err != nil {
		slog.Error(err.Error())
		goto MainMenu
	}

	user, err = handler.ReadUserByEmail(user.Email)
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

	u, err = handler.ReadUserByEmail(user.Email)
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
	fmt.Println("Placeholder")
	return

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
	// case "2":
	// 	goto AddToCart
	// case "3":
	// 	goto ShowCart
	// case "4":
	// 	goto CreateOrder
	default:
		goto Exit
	}

ShowAllProducts:
	fmt.Println("\nShowing all products.....")
	fmt.Fprintln(w, "| Name\t Description\t Price\t")

	products, err := handler.ReadAllProducts()
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

	PrintStdOut(&buf)

	goto UserMenu
}
