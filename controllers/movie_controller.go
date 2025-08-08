package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"test-case/models"
)

var Movies []models.Movie
var currentID = 1

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	var movie models.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	movie.ID = currentID
	currentID++
	Movies = append(Movies, movie)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/movies/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updated models.Movie
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	for i, movie := range Movies {
		if movie.ID == id {
			updated.ID = id
			Movies[i] = updated
			json.NewEncoder(w).Encode(updated)
			return
		}
	}

	http.NotFound(w, r)
}

func ListMovies(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	start := (page - 1) * limit
	end := start + limit

	if start >= len(Movies) {
		json.NewEncoder(w).Encode([]models.Movie{})
		return
	}
	if end > len(Movies) {
		end = len(Movies)
	}

	json.NewEncoder(w).Encode(Movies[start:end])
}


// Checker is contain
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

//Search Movies
func SearchMovies(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		http.Error(w, "Query param 'q' is required", http.StatusBadRequest)
		return
	}

	var results []models.Movie
	for _, movie := range Movies {
		matchTitle := containsIgnoreCase(movie.Title, query)
		matchDescription := containsIgnoreCase(movie.Description, query)
		matchArtist := anyMatch(movie.Artists, query)
		matchGenre := anyMatch(movie.Genres, query)
		matchDuration := strings.Contains(
            strconv.Itoa(movie.Duration),
            query,
        )
		// Log the match variables
		log.Printf(
			"Movie: %s | TitleMatch: %t | DescMatch: %t | ArtistMatch: %t | GenreMatch: %t\n",
			movie.Title, matchTitle, matchDescription, matchArtist, matchGenre,
		)
		if matchTitle || matchDescription || matchArtist || matchGenre || matchDuration {
			results = append(results, movie)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Failed to encode results", http.StatusInternalServerError)
	}
}

// Helper: match any string in slice
func anyMatch(slice []string, query string) bool {
	for _, item := range slice {
		if containsIgnoreCase(item, query) {
			return true
		}
	}
	return false
}

//Search Movies Spesific by Type
func SearchMoviesByType(w http.ResponseWriter, r *http.Request) {
	searchType := strings.TrimSpace(r.URL.Query().Get("search_type"))
	query := strings.TrimSpace(r.URL.Query().Get("q"))

	if searchType == "" || query == "" {
		http.Error(w, "Both 'search_type' and 'q' query params are required", http.StatusBadRequest)
		return
	}

	var results []models.Movie

	for _, movie := range Movies {
		var match bool

		switch strings.ToLower(searchType) {
		case "title":
			match = containsIgnoreCase(movie.Title, query)
		case "description":
			match = containsIgnoreCase(movie.Description, query)
		case "artist":
			match = anyMatch(movie.Artists, query)
		case "genre":
			match = anyMatch(movie.Genres, query)
		default:
			http.Error(w, "Invalid search_type. Must be 'title', 'description', 'artist', or 'genre'", http.StatusBadRequest)
			return
		}

		// Log match result
		log.Printf("Movie: %s | SearchType: %s | Match: %t\n", movie.Title, searchType, match)

		if match {
			results = append(results, movie)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Failed to encode results", http.StatusInternalServerError)
	}
}


//DB
func CreateMovieDB(w http.ResponseWriter, r *http.Request) {
	var movie models.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO movies (title, description, artists, genres, duration) 
			VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := models.DB.QueryRow(query, movie.Title, movie.Description, models.PqArray(movie.Artists), models.PqArray(movie.Genres), movie.Duration).Scan(&movie.ID)
	if err != nil {
		http.Error(w, "Failed to insert movie", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}

func UpdateMovieDB(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/moviesdb/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updated models.Movie
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `UPDATE movies SET title=$1, description=$2, artists=$3, genres=$4, duration=$5 WHERE id=$6`
	_, err = models.DB.Exec(query, updated.Title, updated.Description, models.PqArray(updated.Artists), models.PqArray(updated.Genres), updated.Duration, id)
	if err != nil {
		http.Error(w, "Failed to update movie", http.StatusInternalServerError)
		return
	}

	updated.ID = id
	json.NewEncoder(w).Encode(updated)
}

func ListMoviesDB(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	rows, err := models.DB.Query(`SELECT id, title, description, artists, genres FROM movies LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		http.Error(w, "Failed to fetch movies", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		var artists, genres []string
		err := rows.Scan(&m.ID, &m.Title, &m.Description, models.PqArrayScan(&artists), models.PqArrayScan(&genres))
		if err != nil {
			http.Error(w, "Error scanning movie", http.StatusInternalServerError)
			return
		}
		m.Artists = artists
		m.Genres = genres
		movies = append(movies, m)
	}

	json.NewEncoder(w).Encode(movies)
}

func SearchMoviesDB(w http.ResponseWriter, r *http.Request) {
	queryParam := strings.TrimSpace(r.URL.Query().Get("q"))
	if queryParam == "" {
		http.Error(w, "Query param 'q' is required", http.StatusBadRequest)
		return
	}

	rows, err := models.DB.Query(`
		SELECT id, title, description, artists, genres, duration
		FROM movies 
		WHERE title ILIKE '%' || $1 || '%' 
			OR description ILIKE '%' || $1 || '%' 
			OR ($1 ~ '^[0-9]+$' AND duration = $1::int)
			OR EXISTS (SELECT 1 FROM unnest(artists) a WHERE a ILIKE '%' || $1 || '%')
			OR EXISTS (SELECT 1 FROM unnest(genres) g WHERE g ILIKE '%' || $1 || '%')
	`, queryParam)
	if err != nil {
		http.Error(w, "Failed to search movies", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []models.Movie
	for rows.Next() {
		var m models.Movie
		var artists, genres []string
		err := rows.Scan(&m.ID, &m.Title, &m.Description, models.PqArrayScan(&artists), models.PqArrayScan(&genres), &m.Duration)
		if err != nil {
			http.Error(w, "Error scanning movie", http.StatusInternalServerError)
			return
		}
		log.Printf("Movie: %s | Artists: %v | Genres: %v", m.Title, artists, genres)
		m.Artists = artists
		m.Genres = genres
		results = append(results, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func SearchMoviesByTypeDB(w http.ResponseWriter, r *http.Request) {
	searchType := strings.TrimSpace(r.URL.Query().Get("search_type"))
	queryParam := strings.TrimSpace(r.URL.Query().Get("q"))

	if searchType == "" || queryParam == "" {
		http.Error(w, "Both 'search_type' and 'q' query params are required", http.StatusBadRequest)
		return
	}

	var sqlQuery string
	var rows *sql.Rows
	var err error

	// Build query dynamically based on search_type
	switch strings.ToLower(searchType) {
	case "title":
		sqlQuery = `SELECT id, title, description, artists, genres, duration
		            FROM movies WHERE title ILIKE '%' || $1 || '%'`
		rows, err = models.DB.Query(sqlQuery, queryParam)

	case "description":
		sqlQuery = `SELECT id, title, description, artists, genres, duration
		            FROM movies WHERE description ILIKE '%' || $1 || '%'`
		rows, err = models.DB.Query(sqlQuery, queryParam)

	case "artist":
		sqlQuery = `SELECT id, title, description, artists, genres, duration
		            FROM movies 
		            WHERE EXISTS (SELECT 1 FROM unnest(artists) a WHERE a ILIKE '%' || $1 || '%')`
		rows, err = models.DB.Query(sqlQuery, queryParam)

	case "genre":
		sqlQuery = `SELECT id, title, description, artists, genres, duration
		            FROM movies 
		            WHERE EXISTS (SELECT 1 FROM unnest(genres) g WHERE g ILIKE '%' || $1 || '%')`
		rows, err = models.DB.Query(sqlQuery, queryParam)
	case "duration":
		sqlQuery = `SELECT id, title, description, artists, genres, duration
		            FROM movies 
		            WHERE $1 ~ '^[0-9]+$'
		  			AND duration = CAST($1 AS int)`
		rows, err = models.DB.Query(sqlQuery, queryParam)

	default:
		http.Error(w, "Invalid search_type. Must be 'title', 'description', 'artist', or 'genre'", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Failed to search movies", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []models.Movie
	for rows.Next() {
		var m models.Movie
		var artists, genres []string
		if err := rows.Scan(&m.ID, &m.Title, &m.Description, models.PqArrayScan(&artists), models.PqArrayScan(&genres), &m.Duration); err != nil {
			http.Error(w, "Error scanning movie", http.StatusInternalServerError)
			return
		}
		m.Artists = artists
		m.Genres = genres
		results = append(results, m)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Failed to encode results", http.StatusInternalServerError)
	}
}