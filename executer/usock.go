package executer

import (
	"net"
	"os"
	"strconv"
)

func Listen(sock_path string) (net.Listener, error) {
	_, err := os.Stat(sock_path)
	if err == nil {
		err = os.Remove(sock_path)
		if err != nil {
			return nil, err
		}
	}

	l, err := net.Listen("unix", sock_path)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func Recv(conn net.Conn, content []byte) (int, error) {
	recv_message := ""
	scan := make([]byte, 1)
	for {
		_, err := conn.Read(scan)
		if err != nil {
			return -1, err
		}
		if string(scan) == "\n" {
			break
		}
		recv_message += string(scan)
	}

	contentLength, err := strconv.Atoi(recv_message)
	if err != nil {
		return -1, err
	}

	recv_size := 0
	recv_message = ""
	for {
		n, err := conn.Read(content)
		if err != nil {
			return -1, err
		}

		recv_message += string(content[:n])
		recv_size += n
		if recv_size >= contentLength {
			break
		}
	}

	content = []byte(recv_message)
	return len(recv_message), nil
}

func Send(conn net.Conn, content []byte) (int, error) {
	contentLength := len(content)
	_, err := conn.Write([]byte(strconv.Itoa(contentLength) + "\n"))
	if err != nil {
		return -1, err
	}

	_, err = conn.Write(content)
	if err != nil {
		return -1, err
	}

	return contentLength, nil
}
