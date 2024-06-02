package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/GustasG/waifus/internal/handlers"
)

func main() {
	h, err := handlers.NewPageHandler()
	if err != nil {
		log.Fatalf("could not create page handler: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("GET /", h.HandleIndex())
	mux.Handle("GET /language/{language}", h.HandleLanguage())
	mux.Handle("GET /assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))))
	mux.Handle("GET /favicon.ico", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "assets/favicon.ico")
	}))

	s := &http.Server{
		Addr:         ":5000",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		defer close(idleConnsClosed)

		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt)
		<-shutdown

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.SetKeepAlivesEnabled(false)
		if err := s.Shutdown(ctx); err != nil {
			log.Printf("shutdown error: %v", err)
		}
	}()

	log.Printf("starting server on %s", s.Addr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("server listening error: %v", err)
	}

	<-idleConnsClosed
	log.Println("server shutting down")
}
