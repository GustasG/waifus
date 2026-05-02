package index

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	languages   []string
	counts      map[string]int
	totalImages int
}

func NewHandler(languages []string, counts map[string]int, totalImages int) Handler {
	return Handler{
		languages:   languages,
		counts:      counts,
		totalImages: totalImages,
	}
}

func languagesJSON(languages []string) string {
	b, _ := json.Marshal(languages)
	return string(b)
}

func (h Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Header.Get("Hx-Request") == "true" {
		component := indexContent(h.languages, h.counts, h.totalImages)
		component.Render(r.Context(), w)
	} else {
		component := indexPage(h.languages, h.counts, h.totalImages)
		component.Render(r.Context(), w)
	}
}
