package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/javadmohebbi/goAuth"
	"github.com/sirupsen/logrus"
)

// handle login to server
func (s *Server) LoginController(w http.ResponseWriter, r *http.Request) {
	s.debugger.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)

	// Set header content type to application/json
	w.Header().Set("Content-Type", "application/json")

	type body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var b body
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		// err in decoding body
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"result":  false,
			"message": HTTP_BAD_REQUEST_BODY,
		})
		return
	}

	// prepare user model
	u := goAuth.User{Username: b.Username, Password: b.Password}

	// sign in to system
	t, err := u.SignIn(s.db)

	// check sign in result
	if err != nil {
		// Invalid credentials
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"result":  false,
			"message": INVALID_CREDENTIALS,
		})
		return
	}

	// successful signin
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"result": true,
		"token":  t,
	})
	return

}
