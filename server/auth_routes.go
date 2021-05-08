package server

import "github.com/gorilla/mux"

func (s *Server) AuthRoutes(r *mux.Router) {

	// login
	r.HandleFunc("/login", s.LoginController).Methods("POST", "OPTIONS")
}
