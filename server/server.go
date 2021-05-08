package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/javadmohebbi/goAuth"
	"github.com/javadmohebbi/goAuth/server/debugger"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

var API_VERSION = 1

// goAuht http web server
type Server struct {
	host string
	port int

	debugger *debugger.Debugger

	db *gorm.DB
}

// create new server
func New(host string, port int, debug *debugger.Debugger, db *gorm.DB) *Server {
	debug.Verbose("Create new API server...!", logrus.InfoLevel)
	return &Server{
		host:     host,
		port:     port,
		debugger: debug,
		db:       db,
	}
}

// serve http for our API server
func (s *Server) ServeHTTP() {
	s.debugger.Verbose("Serving HTTP web server", logrus.InfoLevel)

	// new mux router
	mr := mux.NewRouter()

	// new api routes
	apiRoutes := mr.PathPrefix(fmt.Sprintf("/v%d/api", API_VERSION)).Subrouter()

	// auth routes
	authRoutes := apiRoutes.PathPrefix("/auth").Subrouter()
	s.AuthRoutes(authRoutes)

	// 404 | notfound routes shows all api routes in json format
	mr.Path("/").Name("notfound").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.debugger.Verbose(fmt.Sprintf("URI %v called!", r.URL.RequestURI()), logrus.DebugLevel)
		// pinting all routes to console!

		var routeStr []map[string]string
		mr.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			t, err := route.GetPathTemplate()
			if err != nil {
				return err
			}

			methods := ""
			httpmethods, err := route.GetMethods()
			if err == nil {
				for _, m := range httpmethods {
					methods += m + " "
				}
				routeStr = append(routeStr, map[string]string{
					"path":    t,
					"methods": strings.TrimSpace(methods),
					"name":    route.GetName(),
				})
			} else {
				methods = fmt.Sprintf("Error: %v", err)
			}

			return nil
		})
		w.Header().Set("Content-Type", "application/json")
		ret := map[string]interface{}{
			"routes": routeStr,
		}

		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(ret)
		return
	})

	// allow all CORS
	c := cors.AllowAll()

	// log all routes
	mr.Use(loggingMiddleware(s))

	// http server configuration
	srv := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", s.host, s.port),
		Handler: c.Handler(mr),
	}

	var err error
	if isHTTPS, _ := strconv.ParseBool(os.Getenv(goAuth.GOAUTH_CERT_ENABLED)); isHTTPS {
		// TLS
		s.debugger.Verbose(fmt.Sprintf("CERT:'%v'\nKEY:'%v'", os.Getenv(goAuth.GOAUTH_CERT_PATH), os.Getenv(goAuth.GOAUTH_CERT_KEY)), logrus.DebugLevel)
		s.debugger.Verbose(fmt.Sprintf("Listening on HTTPS server %v:%v", s.host, s.port), logrus.InfoLevel)

		err = srv.ListenAndServeTLS(os.Getenv(goAuth.GOAUTH_CERT_PATH), os.Getenv(goAuth.GOAUTH_CERT_KEY))

	} else {
		// NO TLS
		s.debugger.Verbose(fmt.Sprintf("Listening on HTTP server %v:%v", s.host, s.port), logrus.InfoLevel)

		err = srv.ListenAndServe()
	}

	if err != nil {
		s.debugger.Verbose(fmt.Sprintf("Can not serve HTTP(s) due to error: %v", err.Error()), logrus.ErrorLevel)
		os.Exit(1001)
	}

}
