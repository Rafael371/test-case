package main

import (
	"fmt"
	"net/http"
	"test-case/routes"
)

func main() {
	routes.RegisterRoutes()

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}