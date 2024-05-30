package executer

import (
	"fmt"
	"os"
	"os/signal"
)

type Closer struct {
	CloserMap map[string](*Executer)
	Ch        chan os.Signal
}

func (self *Closer) Init() {
	self.CloserMap = make(map[string]*Executer)
}

func (self *Closer) Append(runtime string, exec *Executer) {
	self.CloserMap[runtime] = exec
}

func (self *Closer) Delete(runtime string) {
	os.Remove(self.CloserMap[runtime].Path)
	self.CloserMap[runtime].Listener.Close()
	delete(self.CloserMap, runtime)
}

func (self *Closer) PanicHandler() {
	self.Ch = make(chan os.Signal, 1)
	signal.Notify(self.Ch, os.Interrupt)
	go func() {
		for sig := range self.Ch {
			fmt.Println("signal received!!", sig)

			close(self.Ch)
			for _, exec := range self.CloserMap {
				os.Remove(exec.Path)
				exec.Listener.Close()
			}
			panic("emergency stop!!")
		}
	}()
}
