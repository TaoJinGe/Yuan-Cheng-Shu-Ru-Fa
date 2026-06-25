package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Port                  string
	DBFile                string
	RecordsDir            string
	TokenExpireDays       int
	ShortTokenExpireHours int
}

func Load() Config {
	port := flag.String("port", env("PORT", "8080"), "HTTP listen port")
	dbFile := flag.String("db-file", env("DB_FILE", "data/app.db"), "sqlite database file")
	recordsDir := flag.String("records-dir", env("RECORDS_DIR", "records"), "sent text records directory")
	days := flag.Int("token-days", envInt("TOKEN_EXPIRE_DAYS", 30), "remember token days")
	hours := flag.Int("short-token-hours", envInt("SHORT_TOKEN_EXPIRE_HOURS", 24), "short token hours")
	flag.Parse()

	return Config{
		Port:                  *port,
		DBFile:                *dbFile,
		RecordsDir:            *recordsDir,
		TokenExpireDays:       *days,
		ShortTokenExpireHours: *hours,
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
