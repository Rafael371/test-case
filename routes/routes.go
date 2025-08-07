package routes

import (
	"net/http"
	"test-case/controllers"
)

func RegisterRoutes() {
	http.HandleFunc("/convert", controllers.ConvertHandler)
}