package models

//客户端建立起来的映射信息
type MapInfo struct{
	//客户建立任务时，端监听的端口
	ClientListenPort int
	//客户端与服务端通讯时的Socket端口
	ClientSocketPort int
	//客户端希望服务端暴露的端口
	ServerListenPort int
}