package language

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

	images := make([]string, 0, len(entries))
	for _, entry := range entries {
		fullPath := fmt.Sprintf("/assets/languages/%s/%s", url.PathEscape(language), url.PathEscape(entry.Name()))
		images = append(images, fullPath)
	}

	return images, nil
}

func NewPageHandler() (LanguagePageHandler, error) {
	entries, err := os.ReadDir(filepath.Join("assets", "languages"))
	if err != nil {
		return LanguagePageHandler{}, err
	}

	if len(entries) == 0 {
		return LanguagePageHandler{}, fmt.Errorf("no languages found")
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
			return LanguagePageHandler{}, err
		}
	}

	return LanguagePageHandler{
		languages: languages,
		images:    images,
	}, nil
}

func (h LanguagePageHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("/language/%s", h.languages[0]), http.StatusSeeOther)
}

func (h LanguagePageHandler) HandleLanguage(w http.ResponseWriter, r *http.Request) {
	language := r.PathValue("language")

	images, ok := h.images[language]
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Header.Get("Hx-Request") == "true" {
		component := imageGrid(images)
		component.Render(r.Context(), w)
	} else {
		component := languagePage(h.languages, images, language)
		component.Render(r.Context(), w)
	}
}
