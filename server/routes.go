package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	SuperAdmin = "super_admin"
	Admin      = "admin"
)

func (s *Server) registerRoutes() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/auth", func(r chi.Router) {
		r.With(s.Allow(SuperAdmin, Admin)).Route("/admin", func(r chi.Router) {
			r.Post("/new", s.CreateAdminAccount)
			r.Post("/login", func(w http.ResponseWriter, r *http.Request) {})
			r.Post("/password", func(w http.ResponseWriter, r *http.Request) {})
		})
	})

	r.Route("/hook", func(r chi.Router) {
		r.Get("/", s.Verify)
		r.Post("/", s.Handle)
	})

	r.With(s.Allow(SuperAdmin, Admin)).Route("/categories", func(r chi.Router) {
		r.Post("/", s.handleCategoryCreate)
		r.Get("/", s.handleCategoryGetAll)
	})

	s.Router = r
}
