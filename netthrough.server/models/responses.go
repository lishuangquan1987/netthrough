package models

type RegisterResponse struct{
	ClientId string
	IsSuccess bool
	ErrMsg string
}

type StatusCheckResponse struct {
	ClientId string
	IsSuccess bool
	ErrMsg string
	//是否有外部请求的数据
	HasData bool
	//一次外部请求一个ID,回应的时候按照此ID去回复,如果有多个，则表示外部并发请求
	SessionId []string
}
type ReadDataResponse struct{
	ClientId string
	SessionId string
	IsSuccess bool
	ErrMsg string
	HasData bool
	Data []byte
}
type WriteDataResponse struct{
	ClientId string
	SessionId string
	IsSuccess bool
	ErrMsg string
}
type UnRegisterResponse struct{
	ClientId string
	IsSuccess bool
	ErrMsg string
}