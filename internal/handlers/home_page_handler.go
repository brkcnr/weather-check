package handlers

import (
	"net/http"
	"path/filepath"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	htmlPath := filepath.Join("static", "index.html")
	http.ServeFile(w, r, htmlPath)
}
