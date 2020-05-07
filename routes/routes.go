package routes

import (
	"github.com/vicentedpsantos/news-app-go/handlers"
	"net/http"
)

func StartRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/search", handlers.SearchHandler)
	mux.HandleFunc("/", handlers.IndexHandler)
}
