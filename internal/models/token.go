package models

type GenerateRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
}
