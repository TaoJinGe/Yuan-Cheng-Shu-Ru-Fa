package web

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"

	"voice-bridge-server/internal/auth"
	"voice-bridge-server/internal/config"
	"voice-bridge-server/internal/hub"
)

type Handler struct {
	cfg      config.Config
	sessions *auth.Store
	hub      *hub.Hub
	upgrader websocket.Upgrader
}

func NewHandler(cfg config.Config, sessions *auth.Store, hub *hub.Hub) *Handler {
	return &Handler{
		cfg: cfg, sessions: sessions, hub: hub,
		upgrader: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
	}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.health)
	mux.HandleFunc("/api/register", h.register)
	mux.HandleFunc("/api/login", h.login)
	mux.HandleFunc("/api/logout", h.logout)
	mux.HandleFunc("/ws", h.ws)
	mux.Handle("/", http.FileServer(http.Dir("static")))
	return mux
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("ok"))
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Remember bool   `json:"remember"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	session, ok := h.sessions.Login(req.Username, req.Password, req.Remember)
	if !ok {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}
	writeJSON(w, session)
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Remember bool   `json:"remember"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	session, _, err := h.sessions.Register(req.Username, req.Password, req.Remember)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, session)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Token string `json:"token"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	h.sessions.Logout(req.Token)
	writeJSON(w, map[string]bool{"ok": true})
}

func (h *Handler) ws(w http.ResponseWriter, r *http.Request) {
	ws, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	h.hub.Serve(ws)
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(value)
}
