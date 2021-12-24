package utils

import "net"

func HandleConnection(r, w net.Conn) {
	defer r.Close()
	defer w.Close()

	var buffer = make([]byte, 100000)
	for {
		n, err := r.Read(buffer)
		if err != nil {
			break
		}

		n, err = w.Write(buffer[:n])
		if err != nil {
			break
		}
	}

}
