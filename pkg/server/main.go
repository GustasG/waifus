package main

import (
	"log"
	"net/http"

	"github.com/GustasG/waifus/internal/handlers"
)

func main() {
	h, err := handlers.NewPageHandler()
	if err != nil {
		log.Fatalf("could not create page handler: %v", err)
	}

	http.Handle("/", h.HandleIndex())
	http.Handle("/language/{language}", h.HandleLanguage())
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))

	http.Handle("/favicon.ico", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "assets/favicon.ico")
	}))

	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
