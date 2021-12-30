package business

import (
	"fmt"
	"net"
	"time"

	"netthrough.client/config"
	"netthrough.client/httphelper"
	"netthrough.client/models"
)

type Task struct {
	ClientId string
	conn     net.Conn
	isStop   bool
}

func (t *Task) Start() error {
	//连接端口
	conn, err := net.Dial("tcp", config.Config.SourceAddr)
	if err != nil {
		fmt.Printf("fail to connect to %s,reason:%v\n", config.Config.SourceAddr, err)
		return err
	}
	t.conn = conn
	//向服务端注册
	var registerResponse models.RegisterResponse
	if err := httphelper.PostObj(fmt.Sprintf("http://%s:5001/register", config.Config.ServerIp), models.RegisterRequest{
		ClientId:         t.ClientId,
		ServerListenPort: config.Config.ServerPort,
	}, &registerResponse); err != nil {
		fmt.Printf("register server error,reason:%v\n", err)
		return err
	}
	if !registerResponse.IsSuccess {
		fmt.Printf("register server fail,reason:%v\n", err)
		return err
	}
	go t.process()
	return nil
}
func (t *Task) process() {
	statusCheckRequest := models.StatusCheckRequest{
		ClientId: t.ClientId,
	}
	var statusCheckResponse models.StatusCheckResponse
	for {
		if t.isStop {
			return
		}
		time.Sleep(time.Microsecond * 50)
		//调用服务端StatusCheck接口,看有无数据
		if err := httphelper.PostObj(fmt.Sprintf("http://%s:5001/statuscheck", config.Config.ServerIp), statusCheckRequest, &statusCheckResponse); err != nil {
			fmt.Printf("status check error ,reason:%v\n", err)
			continue
		}
		if !statusCheckResponse.IsSuccess {
			fmt.Printf("status check fail ,reason:%s\n", statusCheckResponse.ErrMsg)
			continue
		}
		if !statusCheckResponse.HasData {
			continue
		}

		for _, sessionId := range statusCheckResponse.SessionId {
			fmt.Printf("had data from server,sessionid:%s\n", sessionId)
		}
		//这里暂时先处理一个请求，并发的后面再考虑
		readDataRequest := models.ReadDataRequest{
			ClientId:  t.ClientId,
			SessionId: statusCheckResponse.SessionId[0],
		}
		var readDataResponse models.ReadDataResponse
		if err := httphelper.PostObj(fmt.Sprintf("http://%s:5001/readdata", config.Config.ServerIp), readDataRequest, &readDataResponse); err != nil {
			fmt.Printf("read data error ,reason:%v\n", err)
			continue
		}
		if !readDataResponse.IsSuccess {
			fmt.Printf("read data fail ,reason:%s\n", readDataResponse.ErrMsg)
			continue
		}
		if !readDataResponse.HasData || len(readDataResponse.Data) == 0 {
			fmt.Printf("read data len is 0 ,reason:%s\n", readDataResponse.ErrMsg)
			continue
		}
		//转发到socket
		if _, err := t.conn.Write(readDataResponse.Data); err != nil {
			fmt.Printf("write data to %s fail,reason:%v\n", t.conn.RemoteAddr(), err)
			continue
		}
		//暂时考虑的是，发送一次就读取一次
		buffer := make([]byte, 1000000)
		if n, err := t.conn.Read(buffer); err != nil {
			fmt.Printf("read data from %s fail ,reason:%v\n", t.conn.RemoteAddr(), err)
			continue
		} else {
			writeDataRequest := models.WriteDataRequest{
				ClientId:  t.ClientId,
				SessionId: statusCheckResponse.SessionId[0],
				Data:      buffer[:n],
			}
			var writeDataResponse models.WriteDataResponse
			err = httphelper.PostObj(fmt.Sprintf("http://%s:5001/writedata", config.Config.ServerIp), writeDataRequest, &writeDataResponse)
			if err != nil {
				fmt.Printf("error to writedata to server ,reason :%v\n", err)
				continue
			}
			if !writeDataResponse.IsSuccess {
				fmt.Printf("fail to writedata to server ,reason :%s\n", writeDataResponse.ErrMsg)
				continue
			}
			fmt.Printf("write data to server ,data len:%d\n", n)
		}

	}
}
func (t *Task) Stop() {
	t.isStop = true
	var response models.UnRegisterResponse
	httphelper.PostObj(fmt.Sprintf("http://%s:5001/unregister", config.Config.ServerIp), models.UnRegisterRequest{
		ClientId: t.ClientId,
	}, &response)
}
