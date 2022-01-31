package http

import (
	"moviebase/moviebase/html"
	"net/http"
)

const MB = 1 << 20

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		html.Render(w, r, "index", nil)
	}
}