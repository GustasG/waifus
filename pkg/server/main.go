package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/GustasG/waifus/internal/language"
)

func main() {
	h, err := language.NewPageHandler()
	if err != nil {
		log.Fatalf("could not create page handler: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h.HandleIndex)
	mux.HandleFunc("GET /language/{language}", h.HandleLanguage)
	mux.Handle("GET /assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))
	mux.Handle("GET /favicon.ico", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "assets/favicon.ico")
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	s := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt)
		<-shutdown

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			log.Printf("graceful shutdown failed: %v", err)
		}

		close(idleConnsClosed)
	}()

	log.Printf("starting server on %s", s.Addr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("cannot start the server: %v", err)
	}

	<-idleConnsClosed
	log.Println("server shutting down")
}
