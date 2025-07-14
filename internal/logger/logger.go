package logger

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/keola-dunn/autolog/internal/random"
)

type Logger struct {
	*slog.Logger

	randomGenerator random.ServiceIface
}

func NewLogger() *Logger {
	return &Logger{
		Logger:          slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		randomGenerator: random.NewService(),
	}
}

func (l *Logger) Error(message string, err error) {
	l.Logger.Error(message, slog.Any("error", err))
}

func (l *Logger) Fatal(message string, err error) {
	l.Logger.Error(message, slog.Any("error", err))
	os.Exit(1)
}

// contextKey is a type used as a key to set values in contexts
type contextKey string

const (
	contextKeyRequestId = contextKey("requestId")
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (l *Logger) RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()
		xReqId := r.Header.Get("x-request-id")

		if strings.TrimSpace(xReqId) == "" {
			// client did not send x-request-id, keep it in the logs
			xReqId, _ = l.randomGenerator.RandomUUID()
		}

		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// attach requestId to context for access in downstream entities
		r = r.WithContext(context.WithValue(r.Context(), contextKeyRequestId, xReqId))

		logEntry := l.Logger.With("path", r.URL.Path, "method", r.Method, "requestId", xReqId)
		logEntry.Info("request received")

		next.ServeHTTP(lrw, r)

		logEntry.Info("request finished!",
			"durationMs", time.Since(start).Milliseconds(),
			"statusCode", lrw.statusCode,
		)

	})
}

func GetRequestId(ctx context.Context) string {
	return ctx.Value(contextKeyRequestId).(string)
}
