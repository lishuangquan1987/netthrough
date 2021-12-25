package models

type RegisterRequest struct {
	ClientSocketPort int
	ServerListenPort int
}

type RegisterResponse struct {
	IsSuccess bool
	ErrMsg    string
}
