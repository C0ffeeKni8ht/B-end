package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"unicode/utf8"

	"lab5/models"
)

//go:embed templates/*.html
var templateFS embed.FS

const (
	maxBodySize = 1 << 20
	maxTextLen  = 300
)

type server struct {
	app *models.App
	tpl *template.Template
	mu  sync.Mutex
}

type entityOps struct {
	name   string
	add    func(r *http.Request) (msg string, err error)
	update func(r *http.Request, id int) (found bool, msg string, err error)
	remove func(id int) (found bool, msg string, err error)
	load   func(id int) (values map[string]string, found bool)
}

type linkRow struct {
	ID                               int
	Element, Parameter, Value, Units string
}

type page struct {
	Section    string
	Error      string
	Entity     string
	ID         int
	Values     map[string]string
	Categories []models.CategoryView
	Elements   []models.ElementView
	Parameters []models.ParametersView
	Links      []linkRow
}

func main() {
	addr := flag.String("addr", "localhost:8080", "адреса, на якій слухає сервер")
	dataDir := flag.String("data", "data", "каталог з CSV-файлами")
	flag.Parse()

	app, err := models.NewApp(*dataDir)
	if err != nil {
		log.Fatal(err)
	}
	tpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	s := &server{app: app, tpl: tpl}

	log.Printf("Сервер запущено: http://%s", *addr)
	log.Fatal(http.ListenAndServe(*addr, s.routes()))
}

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", s.index)
	for _, ops := range []entityOps{s.categoryOps(), s.elementOps(), s.parameterOps(), s.linkOps()} {
		mux.HandleFunc("/"+ops.name, s.handleAdd(ops))
		mux.HandleFunc("/"+ops.name+"/edit", s.handleEdit(ops))
		mux.HandleFunc("/"+ops.name+"/delete", s.handleDelete(ops))
	}
	return mux
}

func (s *server) index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		methodNotAllowed(w, "GET, HEAD")
		return
	}
	s.renderIndex(w, http.StatusOK, "", nil, "")
}

func (s *server) renderIndex(w http.ResponseWriter, status int, section string, values map[string]string, msg string) {
	p := s.newPage()
	p.Section, p.Values, p.Error = section, values, msg
	s.render(w, "index.html", status, p)
}

func (s *server) renderEdit(w http.ResponseWriter, status int, entity string, id int, values map[string]string, msg string) {
	p := s.newPage()
	p.Entity, p.ID, p.Values, p.Error = entity, id, values, msg
	s.render(w, "edit.html", status, p)
}

func (s *server) newPage() page {
	s.mu.Lock()
	defer s.mu.Unlock()

	p := page{
		Categories: s.app.CategoryViews(),
		Elements:   s.app.ElementViews(),
		Parameters: s.app.ParametersViews(),
	}
	paramsByID := make(map[int]models.ParametersView, len(p.Parameters))
	for _, pr := range p.Parameters {
		paramsByID[pr.ID] = pr
	}
	elementNames := make(map[int]string, len(p.Elements))
	for _, e := range p.Elements {
		elementNames[e.ID] = e.Name
	}
	for _, l := range s.app.ElementParametersViews() {
		pr := paramsByID[l.ParameterID]
		p.Links = append(p.Links, linkRow{
			ID: l.ID, Element: elementNames[l.ElementID],
			Parameter: pr.Name, Value: l.Value, Units: pr.Units,
		})
	}
	return p
}

func (s *server) render(w http.ResponseWriter, tmpl string, status int, p page) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := s.tpl.ExecuteTemplate(w, tmpl, p); err != nil {
		log.Printf("шаблон %s: %v", tmpl, err)
	}
}

func (s *server) handleAdd(ops entityOps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead:
			http.Redirect(w, r, "/#"+ops.name, http.StatusFound)

		case http.MethodPost:
			if status, msg := parseBody(w, r); status != 0 {
				s.renderIndex(w, status, ops.name, nil, msg)
				return
			}
			s.mu.Lock()
			msg, err := ops.add(r)
			s.mu.Unlock()

			switch {
			case err != nil:
				serverError(w, r, err)
			case msg != "":
				s.renderIndex(w, http.StatusBadRequest, ops.name, valuesOf(r.PostForm), msg)
			default:
				http.Redirect(w, r, "/#"+ops.name, http.StatusSeeOther)
			}

		default:
			methodNotAllowed(w, "GET, HEAD, POST")
		}
	}
}

func (s *server) handleEdit(ops entityOps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))

		switch r.Method {
		case http.MethodGet, http.MethodHead:
			s.mu.Lock()
			values, found := ops.load(id)
			s.mu.Unlock()
			if !found {
				http.Error(w, "Запис не знайдено", http.StatusNotFound)
				return
			}
			s.renderEdit(w, http.StatusOK, ops.name, id, values, "")

		case http.MethodPost:
			if status, msg := parseBody(w, r); status != 0 {
				s.renderEdit(w, status, ops.name, id, nil, msg)
				return
			}
			s.mu.Lock()
			found, msg, err := ops.update(r, id)
			s.mu.Unlock()

			switch {
			case err != nil:
				serverError(w, r, err)
			case !found:
				http.Error(w, "Запис не знайдено", http.StatusNotFound)
			case msg != "":
				s.renderEdit(w, http.StatusBadRequest, ops.name, id, valuesOf(r.PostForm), msg)
			default:
				http.Redirect(w, r, "/#"+ops.name, http.StatusSeeOther)
			}

		default:
			methodNotAllowed(w, "GET, HEAD, POST")
		}
	}
}

func (s *server) handleDelete(ops entityOps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			methodNotAllowed(w, "POST")
			return
		}
		if status, msg := parseBody(w, r); status != 0 {
			s.renderIndex(w, status, ops.name, nil, msg)
			return
		}
		id, _ := strconv.Atoi(r.PostFormValue("id"))

		s.mu.Lock()
		found, msg, err := ops.remove(id)
		s.mu.Unlock()

		switch {
		case err != nil:
			serverError(w, r, err)
		case !found:
			http.Error(w, "Запис не знайдено", http.StatusNotFound)
		case msg != "":
			s.renderIndex(w, http.StatusConflict, ops.name, nil, msg)
		default:
			http.Redirect(w, r, "/#"+ops.name, http.StatusSeeOther)
		}
	}
}

func parseBody(w http.ResponseWriter, r *http.Request) (status int, msg string) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
	err := r.ParseForm()
	if err == nil {
		return 0, ""
	}
	var tooBig *http.MaxBytesError
	if errors.As(err, &tooBig) {
		return http.StatusRequestEntityTooLarge, "Запит завеликий."
	}
	return http.StatusBadRequest, "Не вдалося прочитати форму."
}

func valuesOf(v url.Values) map[string]string {
	m := make(map[string]string, len(v))
	for k, vals := range v {
		if len(vals) > 0 {
			m[k] = vals[0]
		}
	}
	return m
}

func dataValues(d models.Data) map[string]string {
	m := make(map[string]string, len(d))
	for k, v := range d {
		m[k] = fmt.Sprint(v)
	}
	return m
}

func checkText(label, v string, required bool) string {
	if required && v == "" {
		return label + " обов'язкова."
	}
	if utf8.RuneCountInString(v) > maxTextLen {
		return fmt.Sprintf("%s: максимум %d символів.", label, maxTextLen)
	}
	return ""
}

func methodNotAllowed(w http.ResponseWriter, allow string) {
	w.Header().Set("Allow", allow)
	http.Error(w, "Метод не дозволений", http.StatusMethodNotAllowed)
}

func serverError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%s: %v", r.URL.Path, err)
	http.Error(w, "Внутрішня помилка сервера", http.StatusInternalServerError)
}
