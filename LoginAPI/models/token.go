package models

type Token struct {
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	Session      string `json:"session,omitempty"`
	TokenType    string `json:"TokenType,omitempty"`
	Expires      int64  `json:"expires,omitempty"`
}
