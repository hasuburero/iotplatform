package executer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"time"
)

type Executer struct {
	Pid      int
	Cmd      *exec.Cmd
	Path     string
	Runtime  string
	Listener net.Listener
	Len      int
	Output   []byte
}

type Wait struct {
	err error
}

func (self *Executer) ExecuteWithTimeout(Args []string, timeout int, file []byte) (error, int) {
	var err error
	var code int
	var cnch chan net.Conn = make(chan net.Conn)
	var connection net.Conn
	self.Cmd = exec.Command(Args[0], Args[1:]...)
	end := make(chan error, 1)
	var (
		fd1 io.Writer = &bytes.Buffer{}
		fd2 io.Writer = &bytes.Buffer{}
	)
	go func() {
		go func() {
			//buf, err := self.Cmd.CombinedOutput()
			self.Cmd.Stdout = fd1
			self.Cmd.Stderr = fd2
			err := self.Cmd.Start()
			if err != nil {
				end <- errors.New("start error")
				return
			}
			err = self.Cmd.Wait()
			if err != nil {
				end <- errors.New("wait error")
				return
			}
			end <- err
		}()

		conn, err := self.Listener.Accept()
		cnch <- conn
		if err != nil {
			end <- err
			return
		}

		buf := make([]byte, 1024*1024*100)
		n, err := Recv(conn, buf)
		if err != nil {
			end <- err
			return
		}
		recv := string(buf[:n])
		if recv != "ready" {
			end <- err
			return
		}

		_, err = Send(conn, file)
		if err != nil {
			end <- err
			return
		}

		len, err := Recv(conn, buf)
		if err != nil {
			end <- err
			return
		}

		self.Len = len
		self.Output = buf
		end <- nil
		return
	}()

	connection = <-cnch
	select {
	case err = <-end:
		if err != nil {
			code = 1
			fmt.Println("err has occured with execution")
		} else {
			code = 0
			fmt.Println("no err has occured with execution")
			fmt.Println(fd1)
		}
	case <-time.After(time.Duration(timeout) * time.Second):
		//fmt.Println("timeout")
		err = self.Cmd.Process.Kill()
		if err != nil {
			return err, 0
		}

		err = self.ShowFile()
		if err != nil {
			fmt.Println("file error")
			return err, 0
		}
		err = errors.New("timeout")
		code = 2
	}
	connection.Close()
	return err, code
}

func (self *Executer) ShowFile() error {
	fd1, err := os.Open("sample_result.txt")
	if err != nil {
		return err
	}
	defer fd1.Close()
	buf, err := io.ReadAll(fd1)
	self.Output = buf
	self.Len = len(buf)
	fmt.Println("len = ", self.Len)
	fd2, err := os.Open("sample_node_result.txt")
	if err != nil {
		return err
	}
	defer fd2.Close()
	buf, err = io.ReadAll(fd2)
	fmt.Println(string(buf))
	return nil
}

func (self *Executer) Execute(Args []string) error {
	self.Cmd = exec.Command(Args[0], Args[1:]...)
	err := self.Cmd.Start()
	if err != nil {
		return err
	}
	self.Pid = self.Cmd.Process.Pid
	return nil
}

func (self *Executer) Transit(file []byte) error {
	conn, err := self.Listener.Accept()
	if err != nil {
		return err
	}
	buf := make([]byte, 1024*1024*100)
	n, err := Recv(conn, buf)
	if err != nil {
		return err
	}
	recv := string(buf[:n])
	if recv != "ready" {
		return err
	}

	_, err = Send(conn, file)
	if err != nil {
		return err
	}

	self.Len, err = Recv(conn, buf)
	if err != nil {
		return err
	}

	self.Output = buf
	return err
}

func (self *Executer) Kill() {
	self.Cmd.Process.Kill()
}
