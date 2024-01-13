package main

import (
	"log"
	"net/http"
	"os"

	// this is project path mentioned in go.mod / the package name
	"github.com/DevCodeChan/mongo-golang-restAPI/controllers"
	"github.com/julienschmidt/httprouter"

	"github.com/joho/godotenv"
	// "gopkg.in/mgo.v2"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file: %w", err)
	}

	r := httprouter.New()

	r.GET("/", controllers.ServerStatus)
	r.GET("/user/:id", controllers.GetUser)
	r.POST("/user", controllers.CreateUser)
	r.DELETE("/user/:id", controllers.DeleteUser)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("Set your 'PORT' environment variable")
	}

	address := "localhost:" + PORT
	log.Printf("Server is running on %s\n", address)

	if err := http.ListenAndServe("localhost:5000", r); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
