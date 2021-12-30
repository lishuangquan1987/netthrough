package models

type RegisterRequest struct {
	ClientId         string
	ServerListenPort int
}

type StatusCheckRequest struct {
	ClientId string
}
type ReadDataRequest struct{
	ClientId string
	SessionId string
}
type WriteDataRequest struct{
	ClientId string
	SessionId string
	Data []byte
}

type UnRegisterRequest struct{
	ClientId string
	IsSuccess string
	ErrMsg string
}
