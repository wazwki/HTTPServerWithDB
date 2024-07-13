package main

import (
	"fourthtask/internal/db"
	"fourthtask/internal/handlers"
	"log"
	"net/http"
)

const host string = "localhost:8080"

func main() {
	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/register", handlers.HandlerRegister)
	mux.HandleFunc("/users", handlers.HandlerUser)

	log.Println("Server up with address:", host)
	log.Fatal(http.ListenAndServe(host, mux))
}
