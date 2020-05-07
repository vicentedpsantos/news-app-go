package main

import (
	"github.com/vicentedpsantos/news-app-go/routes"
	"net/http"
	"os"
)

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	routes.StartRoutes(mux)

	http.ListenAndServe(":"+port, mux)
}
