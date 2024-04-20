package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) registerRoutes() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/auth", func(r chi.Router) {
		r.With(s.Authenticate).Route("/admin", func(r chi.Router) {
			r.Post("/new", s.handleCreateAdmin)
			r.Post("/login", s.handleAdminLogin)
			r.Post("/password", s.handleChangePassword)
		})
	})

	r.Route("/hook", func(r chi.Router) {
		r.Get("/", s.Verify)
		r.Post("/", s.Handle)
	})

	r.With(s.Authenticate).Route("/categories", func(r chi.Router) {
		r.Post("/", s.handleCategoryCreate)
		r.Get("/", s.handleCategoryGetAll)
	})

	s.Router = r
}
