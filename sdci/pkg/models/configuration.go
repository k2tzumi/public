package models

// Configuration defines the internal parameters for the application.
type Configuration struct {
	Concurrency map[string]int // map of recipe repo full name to concurrency.
	Recipes     map[string]*Recipe
}
