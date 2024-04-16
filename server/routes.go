package server

import "net/http"

func (s *Server) registerRoutes() {
	mux := http.NewServeMux()

	mux.Handle("POST /v1/auth/admin/new", s.Authenticate()(
		s.handleCreateAdmin(),
	))
	mux.Handle("POST /v1/auth/admin", s.handleAdminLogin())
	mux.Handle("PUT /v1/auth/admin/password", s.Authenticate()(
		s.handleChangePassword(),
	))

	mux.Handle("GET /hook", s.Verify())
	mux.Handle("POST /hook", s.Handle())

	mux.Handle("POST /v1/categories", s.handleCategoryCreate())
	s.Router = mux
}
