package index

import "net/http"

type Handler struct {
	languages       []string
	counts          map[string]int
	totalImages     int
	featuredLanguage string
}

func NewHandler(languages []string, counts map[string]int, totalImages int, featuredLanguage string) Handler {
	return Handler{
		languages:        languages,
		counts:           counts,
		totalImages:      totalImages,
		featuredLanguage: featuredLanguage,
	}
}

func (h Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Header.Get("Hx-Request") == "true" {
		component := indexContent(h.languages, h.counts, h.totalImages, h.featuredLanguage)
		component.Render(r.Context(), w)
	} else {
		component := indexPage(h.languages, h.counts, h.totalImages, h.featuredLanguage)
		component.Render(r.Context(), w)
	}
}
