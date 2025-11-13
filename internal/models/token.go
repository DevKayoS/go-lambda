package models

import "os"

type GenerateTokenRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

var SecretKey = func() []byte {
	key := os.Getenv("JWT_SECRET_KEY")
	if key == "" {
		return []byte("secretKey")
	}
	return []byte(key)
}()
