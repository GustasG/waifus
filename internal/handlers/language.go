package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"

	"github.com/GustasG/waifus/internal/templates"
)

type PageHandler struct {
	Languages []string
}

func NewPageHandler() (*PageHandler, error) {
	entries, err := os.ReadDir(filepath.Join("assets", "languages"))
	if err != nil {
		return nil, err
	}

	languages := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			languages = append(languages, entry.Name())
		}
	}

	if len(languages) == 0 {
		return nil, fmt.Errorf("no languages found")
	}

	return &PageHandler{
		Languages: languages,
	}, nil
}

func (h *PageHandler) HandleLanguage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		language := r.PathValue("language")

		if !slices.Contains(h.Languages, language) {
			http.Error(w, "language not found", http.StatusNotFound)
			return
		}

		entries, err := os.ReadDir(filepath.Join("assets", "languages", language))
		if err != nil {
			http.Error(w, "cannot get images", http.StatusInternalServerError)
			return
		}

		images := make([]string, 0)
		for _, entry := range entries {
			if !entry.IsDir() {
				images = append(images, entry.Name())
			}
		}

		if r.Header.Get("Hx-Request") == "true" {
			component := templates.ImageGrid(images, language)
			component.Render(r.Context(), w)
		} else {
			component := templates.Language(h.Languages, images, language)
			component.Render(r.Context(), w)
		}
	})
}
