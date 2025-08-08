package main

import (
	"log"
	"fmt"
	"net/http"
	"test-case/routes"
	"test-case/models"
)

func main() {
	models.InitDB()
	routes.RegisterRoutes()

	fmt.Println("Server started at http://localhost:8080")
	// Start server once, in main
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}