package main

import (
	"fmt"
	"net"
	"strings"

	"netthrough.server/models"
	"netthrough.server/tasks"

	"github.com/gin-gonic/gin"
)

//tcp 客户端连接
var clientSockets []net.Conn

//外部请求的Socket
var requestSocket net.Conn

//外部请求，将数据放入此通道中
var requestChan chan []byte

//客户端回应的内容，放入此通道，用于返回给外部
var responseChan chan []byte

func main() {
	clientSockets = make([]net.Conn, 0)
	requestChan = make(chan []byte)
	responseChan = make(chan []byte)
	//监听10000端口供外部访问
	// listener, err := net.Listen("tcp", "0.0.0.0:10000")
	// if err != nil {
	// 	fmt.Printf("fail to listen 0.0.0.0:10000,reason:%s", err.Error())
	// 	return
	// }

	// go func(l net.Listener) {
	// 	for {
	// 		con, err := l.Accept()
	// 		if err != nil {
	// 			fmt.Printf("fail to accept ,reason :%s", err.Error())
	// 		} else {
	// 			fmt.Printf("request connect :%s", con.RemoteAddr().String())
	// 			//读取数据
	// 			requestSocket = con
	// 		}

	// 	}
	// }(listener)

	//监听5000端口，供客户端去连接通讯
	clientListener, err := net.Listen("tcp", "0.0.0.0:5000")
	if err != nil {
		fmt.Printf("fail to listen 0.0.0.0:5000, reason :%s", err.Error())
		return
	}
	go func(l net.Listener) {
		for {
			con, err := l.Accept()
			if err != nil {
				fmt.Printf("fail to accept client socket, reason :%s", err.Error())
			} else {
				fmt.Printf("client[%s] connected", con.RemoteAddr().String())
				clientSockets = append(clientSockets, con)
			}
		}
	}(clientListener)

	//开启restful api ,供客户端去注册自己
	r := gin.Default()
	r.POST("/register", Register)
	panic(r.Run("0.0.0.0:5001"))
}

func Register(c *gin.Context) {
	var request models.RegisterRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(200, models.RegisterResponse{
			IsSuccess: false,
			ErrMsg:    "parameter error,please ensure contains fields:ClientSocketPort and ServerListenPort",
		})
		return
	}
	//根据客户端请求：1.监听服务端端口,接收外部请求。2.储存客户端连接的Socket,将外部请求内容转发到这个socket
	//3.持续接收客户端连接的Socket,将信息转发给外部端口的Socket
	var clientSocket net.Conn
	var hasClient bool
	var clientIp string = c.ClientIP()
	fmt.Printf("clientSockets lenth:%d\n", len(clientSockets))
	for _, client := range clientSockets {
		if strings.Split(client.RemoteAddr().String(), ":")[0] == clientIp {
			clientSocket = client
			hasClient = true
			break
		}
	}

	if !hasClient {
		c.JSON(200, models.RegisterResponse{
			IsSuccess: false,
			ErrMsg:    fmt.Sprintf("客户端:%s未与服务器连接", clientIp),
		})
		return
	}
	//监听端口，接受外部的Socket
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", request.ServerListenPort))
	if err != nil {
		c.JSON(200, models.RegisterResponse{
			IsSuccess: false,
			ErrMsg:    fmt.Sprintf("服务器监听%d端口失败，原因:%s", request.ServerListenPort, err.Error()),
		})
	}
	fmt.Printf("建立任务：%s->0.0.0.0:%d", clientIp, request.ServerListenPort)
	task := &tasks.TaskInfo{
		ClientSocket:    clientSocket,
		RequestListener: listener,
	}
	tasks.AddTask(task)
	task.Start()
	c.JSON(200, models.RegisterResponse{
		IsSuccess: true,
		ErrMsg:    "成功",
	})
}
