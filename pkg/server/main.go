package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GustasG/waifus/internal/handlers"
)

func main() {
	h, err := handlers.NewPageHandler()
	if err != nil {
		log.Fatalf("could not create page handler: %v", err)
	}

	http.Handle("/", http.RedirectHandler(fmt.Sprintf("/%s", h.Languages[0]), http.StatusSeeOther))
	http.Handle("/{language}", h.HandleLanguage())
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))

	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
