package config

func (c Config) SavedPassword() string {
	if c.PasswordSecret == "" {
		return ""
	}
	password, err := decryptLocal(c.PasswordSecret)
	if err != nil {
		return ""
	}
	return password
}

func (c *Config) SetSavedPassword(password string, remember bool) {
	c.RememberPassword = remember
	c.PasswordSecret = ""
	if !remember || password == "" {
		return
	}
	secret, err := encryptLocal(password)
	if err == nil {
		c.PasswordSecret = secret
	}
}
