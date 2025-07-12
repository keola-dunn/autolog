package logger

import (
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Logger struct {
	*slog.Logger
}

func NewLogger() *Logger {
	return &Logger{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func (l *Logger) Error(message string, err error) {
	l.Logger.Error(message, slog.Any("error", err))
}

func (l *Logger) Fatal(message string, err error) {
	l.Logger.Error(message, slog.Any("error", err))
	os.Exit(1)
}

func (l *Logger) RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()
		xReqId := r.Header.Get("x-request-id")

		if strings.TrimSpace(xReqId) == "" {
			// client did not send x-request-id, keep it in the logs
			xReqIdUUID, _ := uuid.NewUUID() // google UUID
			xReqId = xReqIdUUID.String()
		}

		logEntry := l.Logger.With("path", r.URL.Path, "method", r.Method, "requestId", xReqId)
		logEntry.Info("request received")

		next.ServeHTTP(w, r)

		logEntry.Info("request finished!",
			"durationMs", time.Since(start),
			"statusCode", r.Response.StatusCode,
		)

	})
}
