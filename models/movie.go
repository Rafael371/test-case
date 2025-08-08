package models

type Movie struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Duration    int      `json:"duration"`
	Artists     []string `json:"artists"`
	Genres      []string `json:"genres"`
}