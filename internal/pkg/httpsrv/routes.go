package httpsrv

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

type Route struct {
	Name    string
	Method  string
	Pattern string
	HFunc   http.Handler
	Queries []string
}

func (s *Server) myRoutes() []Route {
	return []Route{
		{
			Name:    "health",
			Method:  "GET",
			Pattern: "/goapp/health",
			HFunc:   s.handlerWrapper(s.handlerHealth),
		},
		{
			Name:    "websocket",
			Method:  "GET",
			Pattern: "/goapp/ws",
			HFunc:   s.handlerWrapper(s.handlerWebSocket),
		},
		{
			Name:    "home",
			Method:  "GET",
			Pattern: "/goapp",
			HFunc:   s.handlerWrapper(s.handlerHome),
		},
	}
}

// Problem #2: Adding CSRF Token Verification
func CSRFTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract CSRF token from the request header
		csrfToken := r.Header.Get("X-CSRF-Token")

		// Validate the presence of the CSRF token
		if csrfToken == "" {
			http.Error(w, "Missing CSRF token", http.StatusForbidden)
			return
		}

		// Validate the CSRF token (static token)
		validToken := "valid_csrf_token"
		if csrfToken != validToken {
			http.Error(w, "Invalid CSRF token", http.StatusForbidden)
			return
		}

		// If the token is valid, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handlerWrapper(handlerFunc func(http.ResponseWriter, *http.Request)) http.Handler {
	return CSRFTokenMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			r := recover()
			if r != nil {
				s.error(w, http.StatusInternalServerError, fmt.Errorf("%v\n%v", r, string(debug.Stack())))
			}
		}()
		handlerFunc(w, r)
	}))
}
