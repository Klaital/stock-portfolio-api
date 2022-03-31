package server

import (
	"github.com/markbates/goth/gothic"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/klaital/stock-portfolio-api/internal/pkg/config"
)

type Server struct {
	Config *config.Config
}

func New(cfg *config.Config) http.Handler {
	s := Server{Config: cfg}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(cfg.Timeout))

	r.Route("/auth", func(r chi.Router) {
		r.Get("/google/callback", s.googleCallback)
		r.Get("/google", gothic.BeginAuthHandler)
		r.Get("/login", staticPageHandler("views/auth/index.html", nil))
	})

	return r
}

func (s *Server) googleCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.WithError(err).Error("Failed to complete user auth")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t, err := template.ParseFiles("views/auth/success.html")
	if err != nil {
		log.WithError(err).WithField("path", "views/auth/success.html").Error("Error loading template")
		return
	}
	err = t.Execute(w, user)
	if err != nil {
		log.WithError(err).WithField("path", "views/auth/success.html").Error("Error executing template")
	}
}

func staticPageHandler(path string, data interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles(path)
		if err != nil {
			log.WithError(err).WithField("path", path).Error("Error loading template")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = t.Execute(w, false)
		if err != nil {
			log.WithError(err).WithField("path", path).Error("Error executing template")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
