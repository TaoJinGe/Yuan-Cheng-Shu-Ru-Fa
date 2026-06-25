package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Session struct {
	Token     string    `json:"token"`
	UserID    string    `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type API struct {
	base string
}

func NewAPI(wsURL string) API {
	u, err := url.Parse(wsURL)
	if err != nil {
		return API{base: wsURL}
	}
	if u.Scheme == "wss" {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}
	u.Path = ""
	u.RawQuery = ""
	return API{base: strings.TrimRight(u.String(), "/")}
}

func (a API) Login(username, password string, remember bool) (Session, error) {
	body, _ := json.Marshal(authBody(username, password, remember))
	res, err := http.Post(a.base+"/api/login", "application/json", bytes.NewReader(body))
	if err != nil {
		return Session{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return Session{}, fmt.Errorf("登录失败：HTTP %d", res.StatusCode)
	}
	var session Session
	return session, json.NewDecoder(res.Body).Decode(&session)
}

func (a API) Register(username, password string, remember bool) (Session, error) {
	body, _ := json.Marshal(authBody(username, password, remember))
	res, err := http.Post(a.base+"/api/register", "application/json", bytes.NewReader(body))
	if err != nil {
		return Session{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return Session{}, fmt.Errorf("注册失败：HTTP %d", res.StatusCode)
	}
	var session Session
	return session, json.NewDecoder(res.Body).Decode(&session)
}

func (a API) Logout(token string) error {
	body, _ := json.Marshal(map[string]any{"token": token})
	res, err := http.Post(a.base+"/api/logout", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	return res.Body.Close()
}

func authBody(username, password string, remember bool) map[string]any {
	return map[string]any{
		"username": username,
		"password": password,
		"remember": remember,
	}
}
