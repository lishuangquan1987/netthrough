package tasks

import (
	"fmt"
	"net"
)

var taskList []*TaskInfo

func init() {
	taskList = make([]*TaskInfo, 0)
}

type TaskInfo struct {
	//客户端连接的Socket
	ClientSocket net.Conn
	//服务端连接的Socket,用于读取和写入外部请求，每次请求都会变化
	serverSocket    net.Conn
	RequestListener net.Listener
	//外部请求发过来的数据放在这里
	RequestChan chan []byte
	//从客户端发送过来的数据放在这里
	ResponseChan  chan []byte
	IsRuning      bool
	isRequestStop bool
}

//待优化：Read的时候阻塞了，如何Stop掉
func (t *TaskInfo) Start() {
	t.IsRuning = true
	t.isRequestStop = false

	go func(t *TaskInfo) {
		for {
			if t.isRequestStop {
				break
			}
			select {
			case buffer := <-t.RequestChan:
				t.ClientSocket.Write(buffer)
			case buffer := <-t.ResponseChan:
				t.serverSocket.Write(buffer)
			}
		}
	}(t)

	go func() {
		for {
			conn, err := t.RequestListener.Accept()
			if err != nil {
				t.IsRuning = false
				break
			}
			fmt.Printf("[outside connect]:%s\n", conn.RemoteAddr().String())
			t.serverSocket = conn
			//client socket不能关闭
			//接收外部的数据，转发到客户端
			//go transferDataToClient(t, conn)
			//接受客户端的数据，转发到外部
			//go transferDataToOutside(t, conn)

		}
	}()

}

func readDataFromOutside(t *TaskInfo, conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 100000)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			break
		}
		t.RequestChan <- buffer[:n] //放到channel里
	}
}

//数据从外部的请求写入到客户端
func transferDataToClient(t *TaskInfo, con net.Conn) {
	var buffer = make([]byte, 100000)
	for {
		if t.isRequestStop {
			t.IsRuning = false
			break
		}
		n, err := con.Read(buffer)
		if err != nil {
			con.Close()
			fmt.Printf("[outside disconnect]\n", con.RemoteAddr().String())
			//出错了，关闭Socket
			t.IsRuning = false
			break
		}
		fmt.Printf("<transferDataToClient> received %d bytes from [%s]\n ", n, con.RemoteAddr())
		_, err = t.ClientSocket.Write(buffer[:n])
		if err != nil {
			t.IsRuning = false
			break
		}
	}
}

//数据从客户端发送到外部的socket
func transferDataToOutside(t *TaskInfo, con net.Conn) {
	var buffer = make([]byte, 100000)
	for {
		if t.isRequestStop {
			t.IsRuning = false
			break
		}
		n, err := t.ClientSocket.Read(buffer)
		if err != nil {
			t.IsRuning = false
			break
		}
		fmt.Printf("<transferDataToOutside> received %d bytes from [%s]\n ", n, t.ClientSocket.RemoteAddr())
		_, err = con.Write(buffer[:n])
		if err != nil {
			fmt.Printf("[outside disconnect]:%s\n", con.RemoteAddr().String())
			t.IsRuning = false
			con.Close()
			break
		}
	}
}

func (t *TaskInfo) Stop() {
	t.isRequestStop = true
}

func AddTask(task *TaskInfo) {
	taskList = append(taskList, task)
}
