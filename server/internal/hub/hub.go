package hub

import (
	"sync"

	"github.com/gorilla/websocket"

	"voice-bridge-server/internal/auth"
	"voice-bridge-server/internal/records"
)

type Conn struct {
	ws      *websocket.Conn
	writeMu sync.Mutex
	userID  string
	kind    ClientKind
}

type Hub struct {
	sessions *auth.Store
	recorder *records.Recorder
	mu       sync.Mutex
	web      map[string]map[*Conn]bool
	desktop  map[string]map[*Conn]bool
}

func New(sessions *auth.Store, recorder *records.Recorder) *Hub {
	return &Hub{
		sessions: sessions,
		recorder: recorder,
		web:      map[string]map[*Conn]bool{},
		desktop:  map[string]map[*Conn]bool{},
	}
}

func (h *Hub) Serve(ws *websocket.Conn) {
	conn := &Conn{ws: ws}
	defer func() {
		h.unregister(conn)
		ws.Close()
	}()

	for {
		var msg Message
		if err := ws.ReadJSON(&msg); err != nil {
			return
		}
		if !h.handle(conn, msg) {
			return
		}
	}
}

func (h *Hub) handle(conn *Conn, msg Message) bool {
	switch msg.Type {
	case "web_join":
		return h.join(conn, msg, KindWeb)
	case "desktop_join":
		return h.join(conn, msg, KindDesktop)
	case "send_text":
		h.sendText(conn, msg)
		return true
	default:
		write(conn, AckMessage{Type: "ack", OK: false, Error: "unknown message type"})
		return true
	}
}

func (h *Hub) join(conn *Conn, msg Message, kind ClientKind) bool {
	session, ok := h.sessions.Verify(msg.Token)
	if !ok {
		write(conn, AckMessage{Type: "ack", OK: false, Error: "invalid token"})
		return false
	}
	h.unregister(conn)
	conn.userID = session.UserID
	conn.kind = kind

	h.mu.Lock()
	target := h.web
	if kind == KindDesktop {
		target = h.desktop
	}
	if target[conn.userID] == nil {
		target[conn.userID] = map[*Conn]bool{}
	}
	target[conn.userID][conn] = true
	h.mu.Unlock()

	write(conn, AckMessage{Type: "ack", OK: true})
	h.broadcastPresence(conn.userID)
	return true
}

func (h *Hub) sendText(conn *Conn, msg Message) {
	if conn.userID == "" {
		write(conn, AckMessage{Type: "ack", OK: false, Error: "not joined"})
		return
	}
	h.mu.Lock()
	targets := make([]*Conn, 0)
	for desktop := range h.desktop[conn.userID] {
		targets = append(targets, desktop)
	}
	h.mu.Unlock()

	for _, desktop := range targets {
		write(desktop, Message{Type: "insert_text", Text: msg.Text})
	}
	if h.recorder != nil {
		_ = h.recorder.Save(conn.userID, msg.Text)
	}
	write(conn, AckMessage{Type: "ack", OK: true, Delivered: len(targets) > 0})
}

func (h *Hub) unregister(conn *Conn) {
	if conn.userID == "" {
		return
	}
	userID := conn.userID
	h.mu.Lock()
	delete(h.web[userID], conn)
	delete(h.desktop[userID], conn)
	h.mu.Unlock()
	conn.userID = ""
	h.broadcastPresence(userID)
}

func (h *Hub) broadcastPresence(userID string) {
	h.mu.Lock()
	webOnline := len(h.web[userID]) > 0
	desktopOnline := len(h.desktop[userID]) > 0
	targets := make([]*Conn, 0, len(h.web[userID])+len(h.desktop[userID]))
	for conn := range h.web[userID] {
		targets = append(targets, conn)
	}
	for conn := range h.desktop[userID] {
		targets = append(targets, conn)
	}
	h.mu.Unlock()

	status := "disconnected"
	if webOnline && desktopOnline {
		status = "connected"
	}
	msg := PresenceMessage{
		Type: "presence", Status: status,
		WebOnline: webOnline, DesktopOnline: desktopOnline,
	}
	for _, conn := range targets {
		write(conn, msg)
	}
}

func write(conn *Conn, value any) {
	conn.writeMu.Lock()
	defer conn.writeMu.Unlock()
	_ = conn.ws.WriteJSON(value)
}
