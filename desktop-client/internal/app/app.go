package app

import (
	"sync"
	"time"

	"voice-bridge-client/internal/client"
	localconfig "voice-bridge-client/internal/config"
)

type State struct {
	defaultServer string
	cfg           localconfig.Config
	api           client.API
	stop          chan struct{}
	mu            sync.Mutex
}

type LoginInput struct {
	ServerURL        string
	Username         string
	Password         string
	Remember         bool
	RememberPassword bool
}

func New(defaultServer string) (*State, error) {
	cfg, err := localconfig.Load()
	if err != nil {
		return nil, err
	}
	if cfg.ServerURL == "" {
		cfg.ServerURL = defaultServer
	}
	return &State{defaultServer: defaultServer, cfg: cfg}, nil
}

func (s *State) Config() localconfig.Config {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cfg
}

func (s *State) Register(input LoginInput) error {
	return s.auth(input, true)
}

func (s *State) Login(input LoginInput) error {
	return s.auth(input, false)
}

func (s *State) Logout() {
	s.StopSocket()
	s.mu.Lock()
	token := s.cfg.Token
	serverURL := s.cfg.ServerURL
	s.cfg.Token = ""
	s.cfg.UserID = ""
	s.cfg.ExpiresAt = time.Time{}
	_ = localconfig.Save(s.cfg)
	s.mu.Unlock()
	_ = client.NewAPI(serverURL).Logout(token)
}

func (s *State) StartSocket(events client.Events) {
	s.StopSocket()
	cfg := s.Config()
	stop := make(chan struct{})
	s.mu.Lock()
	s.stop = stop
	s.mu.Unlock()
	go func() {
		_ = client.RunDesktopSocket(cfg.ServerURL, cfg.Token, events, stop)
	}()
}

func (s *State) Reconnect(events client.Events) {
	s.StartSocket(events)
}

func (s *State) StopSocket() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.stop != nil {
		close(s.stop)
		s.stop = nil
	}
}

func (s *State) auth(input LoginInput, register bool) error {
	api := client.NewAPI(input.ServerURL)
	var session client.Session
	var err error
	if register {
		session, err = api.Register(input.Username, input.Password, input.Remember)
	} else {
		session, err = api.Login(input.Username, input.Password, input.Remember)
	}
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.cfg.ServerURL = input.ServerURL
	s.cfg.SetSession(session)
	s.cfg.SetSavedPassword(input.Password, input.RememberPassword)
	err = localconfig.Save(s.cfg)
	s.mu.Unlock()
	return err
}
