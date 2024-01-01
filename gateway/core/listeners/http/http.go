package http

import (
	"context"
	"fmt"
	"github.com/rs/cors"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sdq-codes/claimc/gateway/internal/core/errors"
	"github.com/sdq-codes/claimc/gateway/internal/core/logging"
)

const (
	// ErrAddRoutes is the error returned when adding routes to the router fails.
	ErrAddRoutes = errors.Error("failed to add routes")
	// ErrServer is the error returned when the server stops due to an error.
	ErrServer = errors.Error("listen stopped with error")
)

const (
	readHeaderTimeout = 60 * time.Second
)

// Config represents the configuration of the http listener.
type Config struct {
	Port string `yaml:"port"`
}

// Service represents a http service that provides routes for the listener.
type Service interface {
	AddRoutes(r *mux.Router) error
}

// Server represents a http server that listens on a port.
type Server struct {
	server *http.Server
	port   string
}

func contentTypeApplicationJsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
		}
		next.ServeHTTP(w, r)
	})
}

// New instantiates a new instance of Server.
func New(s Service, cfg Config) (*Server, error) {

	r := mux.NewRouter()
	r.Use(contentTypeApplicationJsonMiddleware)
	r.Use(tracingMiddleware)
	r.Use(logTracingMiddleware)
	r.Use(requestLoggingMiddleware)

	if err := s.AddRoutes(r); err != nil {
		return nil, ErrAddRoutes.Wrap(err)
	}

	//c := cors.New(cors.Options{
	//	AllowedOrigins:     []string{"*"},
	//	AllowedHeaders:     []string{"Access-Control-Allow-Origin", "Authorization", "Access-Control-Allow-Origin", "Content-Type"},
	//	AllowCredentials:   true,
	//	AllowedMethods:     []string{http.MethodPut, http.MethodPost, http.MethodGet, http.MethodPatch, http.MethodDelete, http.MethodOptions},
	//	OptionsPassthrough: true,
	//	// Enable Debugging for testing, consider disabling in production
	//	Debug: true,
	//})

	return &Server{
		server: &http.Server{
			Addr: fmt.Sprintf(":%s", cfg.Port),
			BaseContext: func(net.Listener) context.Context {
				baseContext := context.Background()
				return logging.With(baseContext, logging.From(baseContext))
			},
			Handler:           cors.AllowAll().Handler(r),
			ReadHeaderTimeout: readHeaderTimeout,
		},
		port: cfg.Port,
	}, nil
}

// Listen starts the server and listens on the configured port.
func (s *Server) Listen(ctx context.Context) error {
	logging.From(ctx).Info(fmt.Sprintf("http server starting on port: %s", s.port))
	err := s.server.ListenAndServe()
	if err != nil {
		return ErrServer.Wrap(err)
	}

	logging.From(ctx).Info("http server stopped")

	return nil
}
