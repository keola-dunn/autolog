package handler

import "net/http"

type HandlerConfig struct{}

type Handler struct {
}

func NewHandler(cfg HandlerConfig) *Handler {
	return &Handler{}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Autolog API!"))
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy!"))
}
