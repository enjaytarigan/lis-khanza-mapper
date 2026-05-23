package handler

import (
	"encoding/json"
	"net/http"
)

func (s *Server) healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) readyz(w http.ResponseWriter, r *http.Request) {
	if err := s.db.PingContext(r.Context()); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable", "error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeAPI(w http.ResponseWriter, code int, data any, errMsg string) {
	type errBody struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	out := map[string]any{"ok": code >= 200 && code < 300, "data": data, "error": nil}
	if errMsg != "" {
		out["ok"] = false
		out["error"] = errBody{Code: "error", Message: errMsg}
	}
	writeJSON(w, code, out)
}
