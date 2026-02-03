package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
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
	var w *tabwriter.Writer
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

	goto UserMenu

Exit:
	fmt.Println("\nGoodbye!!!!")
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
	w = tabwriter.NewWriter(os.Stdout, 5, 0, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "Name\tDescription\tPrice")
	
	products, err := handler.ReadAllProducts()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	
	for _, product := range products {
		fmt.Fprintf(w, "%s\t%s\t%.2f\n", product.Name, product.Description, product.Price)
	}
	
	if err := w.Flush(); err != nil {
		slog.Error(err.Error())
		goto UserMenu
	}
	
	goto UserMenu
}
