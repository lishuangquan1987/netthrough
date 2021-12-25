package utils

import (
	"fmt"
	"net"
)

func HandleConnection(r, w net.Conn) {
	defer r.Close()
	defer w.Close()

	var buffer = make([]byte, 100000)
	for {
		n, err := r.Read(buffer)
		if err != nil {
			break
		}
		fmt.Printf("received %d bytes from [%s].", n, r.LocalAddr().String())
		n, err = w.Write(buffer[:n])
		if err != nil {
			break
		}
	}

}
