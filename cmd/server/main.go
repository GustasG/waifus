package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/GustasG/waifus/internal/index"
	"github.com/GustasG/waifus/internal/language"
)

func withCacheControl(h http.Handler, maxAge int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
		h.ServeHTTP(w, r)
	})
}

func withHTMLCache(h http.Handler, maxAge int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
		w.Header().Set("Vary", "HX-Request")
		h.ServeHTTP(w, r)
	})
}

func main() {
	langHandler, err := language.NewPageHandler()
	if err != nil {
		log.Fatalf("could not create language handler: %v", err)
	}

	idxHandler := index.NewHandler(langHandler.Languages(), langHandler.Counts(), langHandler.TotalImages())

	mux := http.NewServeMux()
	mux.Handle("GET /", withHTMLCache(http.HandlerFunc(idxHandler.HandleIndex), 60*15))
	mux.Handle("GET /language/{language}", withHTMLCache(http.HandlerFunc(langHandler.HandleLanguage), 60*15))
	mux.Handle("GET /assets/", withCacheControl(
		http.StripPrefix("/assets", http.FileServer(http.Dir("assets"))),
		60*60*24*7,
	))
	mux.Handle("GET /favicon.ico", withCacheControl(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "assets/favicon.ico")
		}),
		60*60*24*7,
	))
	mux.Handle("GET /robots.txt", withCacheControl(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "assets/robots.txt")
		}),
		60*60*24,
	))

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	s := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
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
