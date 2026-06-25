package auth

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type User struct {
	Username     string
	PasswordHash string
}

type UserStore struct {
	db *sql.DB
}

func NewUserStore(path string) (*UserStore, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	store := &UserStore{db: db}
	if err := store.init(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *UserStore) Create(username, passwordHash string) bool {
	_, err := s.db.Exec(
		`insert into users(username, password_hash, created_at) values(?, ?, ?)`,
		username, passwordHash, time.Now().Format(time.RFC3339),
	)
	return err == nil
}

func (s *UserStore) Find(username string) (User, bool) {
	var user User
	err := s.db.QueryRow(
		`select username, password_hash from users where username = ?`,
		username,
	).Scan(&user.Username, &user.PasswordHash)
	return user, err == nil
}

func (s *UserStore) init() error {
	_, err := s.db.Exec(`
		create table if not exists users (
			username text primary key,
			password_hash text not null,
			created_at text not null
		);
		create table if not exists sessions (
			token text primary key,
			user_id text not null,
			expires_at text not null,
			created_at text not null
		);
	`)
	return err
}
