package server

import "github.com/gorilla/mux"

func (s *Server) AuthRoutes(r *mux.Router) {

	// sign in to system and get jwt token from server
	r.HandleFunc("/login", s.LoginController).Methods("POST", "OPTIONS")
}
