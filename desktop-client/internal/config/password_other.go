//go:build !windows

package config

import "fmt"

func encryptLocal(text string) (string, error) {
	_ = text
	return "", fmt.Errorf("local password storage is only supported on Windows")
}

func decryptLocal(secret string) (string, error) {
	_ = secret
	return "", fmt.Errorf("local password storage is only supported on Windows")
}
