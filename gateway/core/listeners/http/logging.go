package http

import (
	"net/http"

	"github.com/sdq-codes/ccasses/gateway/core/logging"
	"go.uber.org/zap"
)

func requestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := logging.WithFields(r.Context(), zap.String("uri", r.RequestURI))

		logging.From(ctx).Info("request", zap.String("method", r.Method))
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
