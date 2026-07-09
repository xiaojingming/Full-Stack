package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/xiaojingming/Full-Stack/server/internal/sql"
)

type envelope struct {
	Data  any    `json:"data"`
	Error string `json:"error,omitempty"`
}

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/health", health)
	mux.HandleFunc("GET /api/system-info", systemInfo)
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic: %v\n%s", rec, debug.Stack())
				writeJSON(w, http.StatusInternalServerError, envelope{Error: "internal server error"})
			}
		}()
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, envelope{Data: map[string]string{"status": "ok"}})
}

func systemInfo(w http.ResponseWriter, r *http.Request) {
	rows, err := sql.Query(r.Context(), "system_info")
	if err != nil {
		writeJSON(w, http.StatusBadGateway, envelope{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, envelope{Data: rows})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("failed to write JSON response: %v", err)
	}
}
