package server

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/javadmohebbi/goAuth"

	"github.com/sirupsen/logrus"
)

// loggingMiddleware will log all of HTTP requests and responses
func loggingMiddleware(a *Server) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		var wr io.Writer
		return handlers.CustomLoggingHandler(wr, next, a.logWriter)
	}
}

// write logs to files
func (a *Server) logWriter(writer io.Writer, params handlers.LogFormatterParams) {

	usr := a.extractUserFromAuthHeader(params.Request)

	adr := params.Request.Header.Get("X-Forwarded-For")
	if adr == "" {
		adr = params.Request.RemoteAddr
	}
	adr, _, err := net.SplitHostPort(adr)
	if err != nil {
		adr = params.Request.RemoteAddr
	}

	// if claims, ok := context.Get(params.Request, "token").(jwt.MapClaims); ok {
	// 	// reqUser.ID = claims["id"].(uint)
	// 	usr = claims["username"].(string)
	// }

	uri := params.Request.RequestURI

	// Requests using the CONNECT method over HTTP/2.0 must use
	// the authority field (aka r.Host) to identify the target.
	// Refer: https://httpwg.github.io/specs/rfc7540.html#CONNECT
	if params.Request.ProtoMajor == 2 && params.Request.Method == "CONNECT" {
		uri = params.Request.Host
	}
	if uri == "" {
		uri = params.URL.RequestURI()
	}

	uri = strings.Split(uri, "?")[0]

	buf := make([]byte, 0, 3*(len(usr)+len(adr)+len(params.Request.Method)+len(uri)+len(params.Request.Proto)+(params.StatusCode)+200)/2)

	// Measurements
	buf = append(buf, "http,"...)

	if usr == "" {
		usr = "NotAuthenticated"
	}

	// Verbose http log
	a.debugger.Verbose("", logrus.InfoLevel, logrus.Fields{
		"user":       usr,
		"client":     adr,
		"uri":        string(appendQuoted([]byte{}, uri)),
		"method":     params.Request.Method,
		"statusCode": strconv.Itoa(params.StatusCode),
		"protocol":   params.Request.Proto,
		"size":       strconv.Itoa(params.Size),
		// "responseTime": fmt.Sprintf("%v", ms),
		"timestamp_ns": fmt.Sprintf("%v", time.Now().UnixNano()),
	})

	// Next line will call customLogWriter but in our case we don't want to use that function
	// writer.Write(buf)
}

func (a *Server) extractUserFromAuthHeader(r *http.Request) string {
	// extract bearer authorization token from request header
	bearerAuthHeader := r.Header.Get("authorization")

	// check if bearer auth header is set
	if bearerAuthHeader == "" {

		// for socket io (web socket)
		tk := r.URL.Query().Get("token")

		if tk == "" {
			return "NotAuthenticated"
		} else {
			// for web socket
			bearerAuthHeader = "Bearer " + tk
		}
	}

	// split bearer and jwt token
	authToken := strings.Split(bearerAuthHeader, " ")
	if len(authToken) != 2 {
		// header set but not in a correct way
		return "NotAuthenticated"
	}

	// check if string first array is Bearer and second is a VALID TOKEN
	if authToken[0] != "Bearer" {
		// header set but not in a correct way
		return "NotAuthenticated"
	}

	// for web socket
	r.Header.Add("Authorization", "Bearer "+fmt.Sprintf("%v", authToken[1]))

	// Check token string
	token, err := jwt.Parse(authToken[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid Token")
		}
		return []byte(os.Getenv(goAuth.GOAUTH_JWT_SECRET_KEY)), nil
	})
	// Check for parsing error
	if err != nil {
		return "NotAuthenticated"
	}
	// validate token if there is not error
	if token.Valid {
		usr, _ := ClaimGetUsername(token.Claims)
		return usr
	}

	return "NotAuthenticated"

}

// append quoted text to logs
const lowerhex = "0123456789abcdef"

// escape URI chars
func appendQuoted(buf []byte, s string) []byte {
	var runeTmp [utf8.UTFMax]byte
	for width := 0; len(s) > 0; s = s[width:] {
		r := rune(s[0])
		width = 1
		if r >= utf8.RuneSelf {
			r, width = utf8.DecodeRuneInString(s)
		}
		if width == 1 && r == utf8.RuneError {
			buf = append(buf, `\x`...)
			buf = append(buf, lowerhex[s[0]>>4])
			buf = append(buf, lowerhex[s[0]&0xF])
			continue
		}
		if r == rune('"') || r == '\\' { // always backslashed
			buf = append(buf, '\\')
			buf = append(buf, byte(r))
			continue
		}
		if strconv.IsPrint(r) {
			n := utf8.EncodeRune(runeTmp[:], r)
			buf = append(buf, runeTmp[:n]...)
			continue
		}
		switch r {
		case '\a':
			buf = append(buf, `\a`...)
		case '\b':
			buf = append(buf, `\b`...)
		case '\f':
			buf = append(buf, `\f`...)
		case '\n':
			buf = append(buf, `\n`...)
		case '\r':
			buf = append(buf, `\r`...)
		case '\t':
			buf = append(buf, `\t`...)
		case '\v':
			buf = append(buf, `\v`...)
		default:
			switch {
			case r < ' ':
				buf = append(buf, `\x`...)
				buf = append(buf, lowerhex[s[0]>>4])
				buf = append(buf, lowerhex[s[0]&0xF])
			case r > utf8.MaxRune:
				r = 0xFFFD
				fallthrough
			case r < 0x10000:
				buf = append(buf, `\u`...)
				for s := 12; s >= 0; s -= 4 {
					buf = append(buf, lowerhex[r>>uint(s)&0xF])
				}
			default:
				buf = append(buf, `\U`...)
				for s := 28; s >= 0; s -= 4 {
					buf = append(buf, lowerhex[r>>uint(s)&0xF])
				}
			}
		}
	}
	return buf

}
