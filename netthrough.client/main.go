package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"netthrough.client/config"
	"netthrough.client/httphelper"
	"netthrough.client/models"
	"netthrough.client/utils"
)

func main() {
	//1.与服务器建立连接,服务器TCP端口：5000
	servercon, err := net.Dial("tcp", fmt.Sprintf("%s:5000", config.Config.ServerIp))
	if err != nil {
		fmt.Printf("fail to connect server.please ensure that netthrough.server is running on server")
		return
	}
	defer servercon.Close()
	//2.连接要转发的地址，得到连接
	sourcecon, err := net.Dial("tcp", config.Config.SourceAddr)
	if err != nil {
		fmt.Printf("fail to connnect %s", config.Config.SourceAddr)
		return
	}
	//2.向服务器注册Socket的远程端口
	clientPort, _ := strconv.ParseInt(strings.Split(servercon.LocalAddr().String(), ":")[1], 10, 64)
	response := models.RegisterResponse{}
	err = httphelper.PostObj(fmt.Sprintf("http://%s:5001/register", config.Config.ServerIp), models.RegisterRequest{
		ClientSocketPort: int(clientPort),
		ServerListenPort: config.Config.ServerPort,
	}, &response)
	if err != nil {
		fmt.Printf("请求出错，原因：%v", err)
		return
	}
	if !response.IsSuccess {
		fmt.Println(response.ErrMsg)
		return
	}
	//处理端口转发
	fmt.Printf("%s<->%s:%d", config.Config.SourceAddr, config.Config.ServerIp, config.Config.ServerPort)
	go utils.HandleConnection(servercon, sourcecon)
	go utils.HandleConnection(sourcecon, servercon)
	//阻塞程序
	select {}

}
