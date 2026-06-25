package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"voice-bridge-client/internal/paste"
)

var ErrAuthExpired = errors.New("登录已失效，请重新登录")

type Events struct {
	Status     func(string)
	Connected  func(bool)
	Inserted   func()
	AuthFailed func()
}

type wsMessage struct {
	Type  string `json:"type"`
	Token string `json:"token,omitempty"`
	Text  string `json:"text,omitempty"`
}

type presenceMessage struct {
	Type          string `json:"type"`
	Status        string `json:"status"`
	WebOnline     bool   `json:"web_online"`
	DesktopOnline bool   `json:"desktop_online"`
}

func RunDesktopSocket(serverURL, token string, events Events, stop <-chan struct{}) error {
	for {
		select {
		case <-stop:
			return nil
		default:
		}
		if err := runOnce(serverURL, token, events, stop); err != nil {
			if errors.Is(err, ErrAuthExpired) {
				call(events.AuthFailed)
				return err
			}
			callStatus(events, fmt.Sprintf("服务器断开：%v", err))
			if waitStop(stop, 3*time.Second) {
				return nil
			}
		}
	}
}

func runOnce(serverURL, token string, events Events, stop <-chan struct{}) error {
	ws, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	if err := ws.WriteJSON(wsMessage{Type: "desktop_join", Token: token}); err != nil {
		return err
	}
	callStatus(events, "已连接服务器，等待手机端")

	for {
		select {
		case <-stop:
			return nil
		default:
		}
		var raw map[string]any
		if err := ws.ReadJSON(&raw); err != nil {
			return err
		}
		switch raw["type"] {
		case "presence":
			handlePresence(raw, events)
		case "insert_text":
			text, _ := raw["text"].(string)
			if err := paste.PasteText(text, 1500*time.Millisecond); err != nil {
				callStatus(events, fmt.Sprintf("自动粘贴失败：%v", err))
			} else {
				callStatus(events, "已粘贴收到的文字")
				call(events.Inserted)
			}
		case "ack":
			if raw["ok"] == false {
				if reason, _ := raw["error"].(string); strings.Contains(reason, "token") {
					return ErrAuthExpired
				}
				return ErrAuthExpired
			}
		}
	}
}

func handlePresence(raw map[string]any, events Events) {
	data, _ := json.Marshal(raw)
	var msg presenceMessage
	_ = json.Unmarshal(data, &msg)
	if msg.Status == "connected" {
		callStatus(events, "绿灯 / 已连接")
		if events.Connected != nil {
			events.Connected(true)
		}
		return
	}
	callStatus(events, "红灯 / 未连接")
	if events.Connected != nil {
		events.Connected(false)
	}
}

func waitStop(stop <-chan struct{}, delay time.Duration) bool {
	select {
	case <-stop:
		return true
	case <-time.After(delay):
		return false
	}
}

func callStatus(events Events, text string) {
	if events.Status != nil {
		events.Status(text)
	}
}

func call(fn func()) {
	if fn != nil {
		fn()
	}
}
