package hub

type ClientKind string

const (
	KindWeb     ClientKind = "web"
	KindDesktop ClientKind = "desktop"
)

type Message struct {
	Type   string `json:"type"`
	RoomID string `json:"roomId,omitempty"`
	Secret string `json:"secret,omitempty"`
	Token  string `json:"token,omitempty"`
	Text   string `json:"text,omitempty"`
}

type PresenceMessage struct {
	Type          string `json:"type"`
	Status        string `json:"status"`
	WebOnline     bool   `json:"web_online"`
	DesktopOnline bool   `json:"desktop_online"`
}

type AckMessage struct {
	Type        string `json:"type"`
	MessageType string `json:"message_type,omitempty"`
	OK          bool   `json:"ok"`
	Delivered   *bool  `json:"delivered,omitempty"`
	Error       string `json:"error,omitempty"`
}
