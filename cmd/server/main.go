package main

import (
	"github.com/gorilla/sessions"
	"github.com/klaital/stock-portfolio-api/internal/pkg/config"
	"github.com/klaital/stock-portfolio-api/internal/pkg/server"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()
	if err := cfg.Validate(); err != nil {
		log.WithError(err).Fatal("Invalid configuration")
	}

	store := sessions.NewCookieStore([]byte(cfg.CookieStoreSecret))
	store.MaxAge(cfg.CookieMaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false // we use HTTP termination at the load balancer

	gothic.Store = store

	goth.UseProviders(
		google.New(cfg.GoogleClientID, cfg.GoogleClientSecret, "http://localhost:3000/auth/google/callback", "email", "profile"),
	)

	handler := server.New(cfg)

	//http.HandleFunc("/auth/google/callback", func(res http.ResponseWriter, req *http.Request) {
	//	if req.Method != "GET" {
	//		res.WriteHeader(http.StatusMethodNotAllowed)
	//		return
	//	}
	//	user, err := gothic.CompleteUserAuth(res, req)
	//	if err != nil {
	//		log.WithError(err).Error("Failed to complete user auth")
	//		return
	//	}
	//	t, err := template.ParseFiles("views/auth/success.html")
	//	if err != nil {
	//		log.WithError(err).WithField("path", "views/auth/success.html").Error("Error loading template")
	//		return
	//	}
	//	err = t.Execute(res, user)
	//	if err != nil {
	//		log.WithError(err).WithField("path", "views/auth/success.html").Error("Error executing template")
	//	}
	//})
	//
	//http.HandleFunc("/auth/google", func(res http.ResponseWriter, req *http.Request) {
	//	gothic.BeginAuthHandler(res, req)
	//})
	//
	//http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
	//	t, err := template.ParseFiles("views/auth/index.html")
	//	if err != nil {
	//		log.WithError(err).WithField("path", "views/auth/index.html").Error("Error loading template")
	//		return
	//	}
	//	err = t.Execute(res, false)
	//	if err != nil {
	//		log.WithError(err).WithField("path", "views/auth/index.html").Error("Error executing template")
	//	}
	//})

	http.ListenAndServe(":3000", handler)
}
