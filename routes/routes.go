package routes

import (
	"net/http"

	"test-case/controllers"
)

func RegisterRoutes() {
	http.HandleFunc("/convert", controllers.ConvertHandler)

	//movies offline (using array)
	http.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			controllers.CreateMovie(w, r)
		} else {
			controllers.ListMovies(w, r)
		}
	})

	http.HandleFunc("/movies/", controllers.UpdateMovie)
	http.HandleFunc("/search", controllers.SearchMovies)
	http.HandleFunc("/search/type", controllers.SearchMoviesByType)

	//movies online (using db)
	http.HandleFunc("/moviesdb", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			controllers.CreateMovieDB(w, r)
		} else {
			controllers.ListMoviesDB(w, r)
		}
	})

	http.HandleFunc("/moviesdb/", controllers.UpdateMovieDB)
	http.HandleFunc("/searchdb", controllers.SearchMoviesDB)
	http.HandleFunc("/searchdb/type", controllers.SearchMoviesByTypeDB)
}