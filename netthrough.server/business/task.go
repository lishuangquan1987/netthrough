package business

import (
	"fmt"
	"net"

	goutils "github.com/typa01/go-utils"
)

var Tasks []*Task

func init() {
	Tasks = make([]*Task, 0)
}

type Task struct {
	ClientId         string
	ServerListenPort int
	Listener         net.Listener
	WriteBuffer      map[string]chan []byte
	ReadBuffer       map[string]chan []byte
}

func (t *Task) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", t.ServerListenPort))
	if err != nil {
		return err
	}

	t.Listener = listener
	t.WriteBuffer = make(map[string]chan []byte)
	t.ReadBuffer = make(map[string]chan []byte)

	go func(task *Task) {
		for {
			conn, err := task.Listener.Accept()
			if err != nil {
				continue
			}
			handleData(t, conn)
		}
	}(t)
	return nil
}
func handleData(t *Task, conn net.Conn) {

	guid := goutils.GUID()

	t.WriteBuffer[guid] = make(chan []byte, 5)
	t.ReadBuffer[guid] = make(chan []byte, 5)

	//读取外部的请求，写入到客户端
	go readOutsideData(t, conn, guid)
	//读取客户端的请求，写入到外部
	go writeDataToOutside(t, conn, guid)
}
func readOutsideData(task *Task, con net.Conn, guid string) {
	bytes := make([]byte, 100000)
	defer delete(task.WriteBuffer, guid)
	defer con.Close()
	for {
		n, err := con.Read(bytes)
		if err != nil {
			break
		}
		fmt.Printf("readOutsideData:<-[%s],len:%d\n", con.RemoteAddr(), n)
		//读取到外部请求的数据，生成一个RequestId
		task.WriteBuffer[guid] <- bytes[:n]
	}
}
func writeDataToOutside(task *Task, con net.Conn, guid string) {
	defer delete(task.ReadBuffer, guid)
	defer con.Close()
	for {
		select {
		case buffer := <-task.ReadBuffer[guid]:
			n, err := con.Write(buffer)
			if err != nil {
				return
			}
			fmt.Printf("writeDataToOutside:->[%s],len:%d\n", con.RemoteAddr(), n)
		}
	}
}

func (t *Task) Stop() error {
	if err := t.Listener.Close(); err != nil {
		return err
	}
	Tasks = delItem(Tasks, t)
	return nil
}

func AddTask(t *Task) {
	Tasks = append(Tasks, t)
}

//从数组中移除元素
func delItem(vs []*Task, s *Task) []*Task {
	for i := 0; i < len(vs); i++ {
		if s == vs[i] {
			vs = append(vs[:i], vs[i+1:]...)
			i = i - 1
		}
	}
	return vs
}
