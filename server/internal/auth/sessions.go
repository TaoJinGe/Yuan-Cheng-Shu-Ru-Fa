package auth

import "time"

func (s *UserStore) SaveSession(session Session) error {
	_, err := s.db.Exec(
		`insert or replace into sessions(token, user_id, expires_at, created_at) values(?, ?, ?, ?)`,
		session.Token,
		session.UserID,
		session.ExpiresAt.Format(time.RFC3339),
		time.Now().Format(time.RFC3339),
	)
	return err
}

func (s *UserStore) FindSession(token string) (Session, bool) {
	var session Session
	var expires string
	err := s.db.QueryRow(
		`select token, user_id, expires_at from sessions where token = ?`,
		token,
	).Scan(&session.Token, &session.UserID, &expires)
	if err != nil {
		return Session{}, false
	}
	expiresAt, err := time.Parse(time.RFC3339, expires)
	if err != nil {
		return Session{}, false
	}
	session.ExpiresAt = expiresAt
	return session, true
}

func (s *UserStore) DeleteSession(token string) {
	_, _ = s.db.Exec(`delete from sessions where token = ?`, token)
}

func (s *UserStore) DeleteExpiredSessions(now time.Time) {
	_, _ = s.db.Exec(`delete from sessions where expires_at <= ?`, now.Format(time.RFC3339))
}
