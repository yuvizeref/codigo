package main

import (
	"codigo/db"
	"codigo/routes"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = db.InitMongoDB()
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v", err)
	}

	r := routes.InitRoutes()

	serverPort := os.Getenv("SERVER_PORT")

	log.Printf("Server is running on http://localhost:%s", serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, r))
}
