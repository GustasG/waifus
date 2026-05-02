package language

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
)

const cdnBase = "https://raw.githubusercontent.com/cat-milk/Anime-Girls-Holding-Programming-Books/master"

type LanguagePageHandler struct {
	languages []string
	images    map[string][]string
	counts    map[string]int
}

func NewPageHandler() (LanguagePageHandler, error) {
	data, err := os.ReadFile(filepath.Join("assets", "manifest.json"))
	if err != nil {
		return LanguagePageHandler{}, fmt.Errorf("manifest not found — run scripts/generate_manifest.py first: %w", err)
	}

	var raw map[string][]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return LanguagePageHandler{}, fmt.Errorf("parse manifest: %w", err)
	}
	if len(raw) == 0 {
		return LanguagePageHandler{}, fmt.Errorf("manifest is empty")
	}

	languages := make([]string, 0, len(raw))
	for lang := range raw {
		languages = append(languages, lang)
	}
	sort.Strings(languages)

	images := make(map[string][]string, len(raw))
	counts := make(map[string]int, len(raw))
	for lang, files := range raw {
		urls := make([]string, len(files))
		for i, f := range files {
			urls[i] = fmt.Sprintf("%s/%s/%s", cdnBase, url.PathEscape(lang), url.PathEscape(f))
		}
		images[lang] = urls
		counts[lang] = len(urls)
	}

	return LanguagePageHandler{languages: languages, images: images, counts: counts}, nil
}

func (h LanguagePageHandler) Languages() []string     { return h.languages }
func (h LanguagePageHandler) Counts() map[string]int  { return h.counts }

func (h LanguagePageHandler) TotalImages() int {
	total := 0
	for _, c := range h.counts {
		total += c
	}
	return total
}

func (h LanguagePageHandler) HandleLanguage(w http.ResponseWriter, r *http.Request) {
	language := r.PathValue("language")

	images, ok := h.images[language]
	if !ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusNotFound)
		notFoundPage(h.languages, h.counts).Render(r.Context(), w)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Header.Get("Hx-Request") == "true" {
		component := imageGrid(images, language, h.counts[language])
		component.Render(r.Context(), w)
	} else {
		component := languagePage(h.languages, h.counts, images, language)
		component.Render(r.Context(), w)
	}
}
