package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"

	"lab5/models"
)

func main() {
	addr := flag.String("addr", "localhost:8000", "адреса, на якій слухає сервер")
	dataDir := flag.String("data", "data", "каталог з CSV-файлами")
	webDir := flag.String("web", "web", "каталог зі статичними файлами (demo.html)")
	flag.Parse()

	app, err := models.NewApp(*dataDir)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/elements", jsonHandler(app.ElementViews, app.ElementView))
	mux.HandleFunc("/api/parameters", jsonHandler(app.ParametersViews, app.ParametersView))
	mux.HandleFunc("/api/element-parameters", jsonHandler(app.ElementParametersViews, app.ElementParametersView))
	mux.Handle("/", http.FileServer(http.Dir(*webDir)))

	log.Printf("Сервер запущено: http://%s", *addr)
	log.Printf("  JSON: /api/elements, /api/parameters, /api/element-parameters")
	log.Printf("  Демо: /demo.html")
	log.Fatal(http.ListenAndServe(*addr, mux))
}

func jsonHandler[T any](list func() []T, one func(id int) (T, bool)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Дозволено лише GET"})
			return
		}

		query := r.URL.Query()
		if !query.Has("id") {
			writeJSON(w, http.StatusOK, list())
			return
		}

		id, err := strconv.Atoi(query.Get("id"))
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id має бути цілим числом"})
			return
		}
		item, ok := one(id)
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Запис не знайдено"})
			return
		}
		writeJSON(w, http.StatusOK, item)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		log.Printf("помилка запису JSON: %v", err)
	}
}
