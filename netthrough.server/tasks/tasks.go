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
	ClientSocket    net.Conn
	RequestListener net.Listener
	IsRuning        bool
	isRequestStop   bool
}

//待优化：Read的时候阻塞了，如何Stop掉
func (t *TaskInfo) Start() {
	t.IsRuning = true
	t.isRequestStop = false

	go func() {
		for {
			conn, err := t.RequestListener.Accept()
			if err != nil {
				t.IsRuning = false
				break
			}
			fmt.Printf("[outside connect]:%s\n", conn.RemoteAddr().String())
			//client socket不能关闭
			//接收外部的数据，转发到客户端
			//go transferDataToClient(t, conn)
			//接受客户端的数据，转发到外部
			//go transferDataToOutside(t, conn)
			go transferData(t, conn)
		}
	}()

}

func transferData(t *TaskInfo, con net.Conn) {
	//从外部接受数据
	var buffer = make([]byte, 100000)
	defer con.Close()
	n, err := con.Read(buffer)
	if err != nil {
		return
	}
	fmt.Printf("<transferData> received %d bytes from [%s]\n ", n, con.RemoteAddr())
	//写入到客户端
	n, err = t.ClientSocket.Write(buffer[:n])
	if err != nil {
		return
	}
	//从客户端读取内容
	n, err = t.ClientSocket.Read(buffer)
	if err != nil {
		return
	}
	//写入到外部
	n, err = con.Write(buffer[:n])
}

//数据从外部的请求写入到客户端
func transferDataToClient(t *TaskInfo, con net.Conn) {

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
