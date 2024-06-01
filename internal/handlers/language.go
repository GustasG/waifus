package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/GustasG/waifus/internal/templates"
)

type LanguagePageHandler struct {
	languages []string
	images    map[string][]string
}

func findImages(language string) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join("assets", "languages", language))
	if err != nil {
		return nil, err
	}

	images := make([]string, 0)
	for _, entry := range entries {
		fullPath := fmt.Sprintf("/assets/languages/%s/%s", url.PathEscape(language), url.PathEscape(entry.Name()))
		images = append(images, fullPath)
	}

	return images, nil
}

func NewPageHandler() (*LanguagePageHandler, error) {
	entries, err := os.ReadDir(filepath.Join("assets", "languages"))
	if err != nil {
		return nil, err
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("no languages found")
	}

	languages := make([]string, 0, len(entries))
	images := make(map[string][]string, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		language := entry.Name()
		languages = append(languages, language)
		images[language], err = findImages(language)
		if err != nil {
			return nil, err
		}
	}

	return &LanguagePageHandler{
		languages: languages,
		images:    images,
	}, nil
}

func (h *LanguagePageHandler) HandleIndex() http.Handler {
	return http.RedirectHandler(fmt.Sprintf("/language/%s", h.languages[0]), http.StatusSeeOther)
}

func (h *LanguagePageHandler) HandleLanguage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		language := r.PathValue("language")

		images, ok := h.images[language]
		if !ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if r.Header.Get("Hx-Request") == "true" {
			component := templates.ImageGrid(images)
			component.Render(r.Context(), w)
		} else {
			component := templates.LanguagePage(h.languages, images, language)
			component.Render(r.Context(), w)
		}
	})
}
