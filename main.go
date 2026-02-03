package main

import (
	"database/sql"
	"log"
	"os"

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
}
