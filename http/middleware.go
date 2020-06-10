package http

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/thetreep/covidtracker"
	"github.com/thetreep/covidtracker/logger"
)

type adapter func(http.Handler) http.Handler

func adapt(h http.Handler, adapters ...adapter) http.Handler {
	// reverse order to apply adapters in order they are specified in parameters
	// (the first adapter is the first one that need to be executed)
	for i := len(adapters) - 1; i >= 0; i-- {
		adapter := adapters[i]
		h = adapter(h)
	}
	return h
}

func (s *Server) routing() adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for path, handler := range s.Routing {
				if path == r.URL.Path {
					handler.ServeHTTP(w, r)
					return
				}
			}
			// h.ServeHTTP(w, r)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{}` + "\n"))
		})
	}
}

func (s *Server) requestId() adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Use existing request ID if available
			requestID := r.Header.Get("X-Request-Id")

			if requestID == "" {
				requestID = uuid.New().String()
			}
			r.WithContext(context.WithValue(r.Context(), logger.ContextKeyRequestId, requestID))
			r.WithContext(context.WithValue(r.Context(), logger.ContextKeyRequestURI, r.RequestURI))
			w.Header().Set("X-Request-Id", requestID)
			h.ServeHTTP(w, r)
		})
	}
}

func (s *Server) ping() adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/ping" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

func (s *Server) cors() adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			addCorsHeader(w)
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			} else {
				h.ServeHTTP(w, r)
			}
		})
	}
}

func (s *Server) auth() adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			tokenStr := r.Header.Get("api-secret")
			if tokenStr == "" {
				Error(r.Context(), w, covidtracker.ErrMissingAPISecret, http.StatusUnauthorized)
				return
			}

			if tokenStr != os.Getenv("THETREEP_COVIDTRACKER_SECRET") {
				Error(r.Context(), w, covidtracker.ErrInvalidAPISecret, http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

func (s *Server) log() adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()
			h.ServeHTTP(w, r)
			path := r.URL.Path
			if _, ignored := pathsIgnoreLogging[path]; ignored {
				return
			}
			stop := time.Since(start)
			latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))

			clientIP := clientIP(r)
			clientUserAgent := r.UserAgent()
			referer := r.Referer()
			hostname, _ := os.Hostname()

			fields := map[string]interface{}{
				"hostname":  hostname,
				"latency":   latency,
				"clientIP":  clientIP,
				"method":    r.Method,
				"path":      path,
				"referer":   referer,
				"userAgent": clientUserAgent,
				"tags":      "http",
			}

			b, err := ioutil.ReadAll(r.Body)
			if err == nil {
				fields["body"] = string(b)
			}

			ctx := r.Context()
			msg := fmt.Sprintf("%s - %s [%s] \"%s %s\" \"%s\" \"%s\" (%dms)", clientIP, hostname, time.Now().Format("02/Jan/2006:15:04:05 -0700"), r.Method, path, referer, clientUserAgent, latency)
			s.logger.InfoWithFields(ctx, fields, msg)
		})
	}
}

func addCorsHeader(w http.ResponseWriter) {
	headers := w.Header()
	headers.Add("Access-Control-Allow-Origin", "*")
	headers.Add("Access-Control-Max-Age", "86400")
	headers.Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	headers.Add("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Access-Control-Allow-Methods, Api-Secret, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
	headers.Add("Access-Control-Expose-Headers", "Content-Length")
	headers.Add("Access-Control-Allow-Credentials", "true")
}

var pathsIgnoreLogging = map[string]struct{}{
	"/":     struct{}{},
	"/ping": struct{}{},
}

func clientIP(r *http.Request) string {
	clientIP := r.Header.Get("X-Forwarded-For")
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	}
	if clientIP != "" {
		return clientIP
	}

	if addr := r.Header.Get("X-Appengine-Remote-Addr"); addr != "" {
		return addr
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

// trim token type "Bearer" if present
func cleanAuthToken(t string) string {
	prefix := "Bearer"
	if len(t) >= len(prefix) && strings.EqualFold(t[0:len(prefix)], prefix) {
		t = t[len(prefix):]
	}
	return strings.TrimSpace(t)
}
