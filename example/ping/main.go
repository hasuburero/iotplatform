package main

import (
	"fmt"
	"github.com/hasuburero/mecrm/api"
	"strconv"
	"strings"
	"time"
)

const (
	num = 10000
	url = "https://mecrm.dolylab.cc/api/v0.5"
)

type Timestamp struct {
	Timer     time.Time
	TimeStamp string
	Hour      int
	Min       int
	Sec       int
	Millis    int
	Ms        int
}

func (t *Timestamp) getms() int {
	t.Ms = t.Millis + 1000*t.Sec + 60000*t.Min + 3600000*t.Hour
	return t.Ms
}

func (t *Timestamp) getTime() {
	t.Timer = time.Now()
	t.TimeStamp = t.Timer.Format(time.StampMilli)
	buf := strings.Split(t.TimeStamp, " ")
	for i, ctx := range buf {
		if ctx == "" {
			buf = append(buf[:i], buf[i+1:]...)
		}
	}
	tbuf := buf[2]
	buf = strings.Split(tbuf, ":")
	t.Hour, _ = strconv.Atoi(string(buf[0]))
	t.Min, _ = strconv.Atoi(string(buf[1]))
	tbuf = buf[2]
	buf = strings.Split(tbuf, ".")
	t.Sec, _ = strconv.Atoi(string(buf[0]))
	t.Millis, _ = strconv.Atoi(string(buf[1]))
}

var logch chan int
var ch [num]chan bool
var wait chan bool

func Log() {
	for _ = range num {
		ms := <-logch
		fmt.Println(ms)
	}
	wait <- true
}
func main() {
	for i := range num {
		ch[i] = make(chan bool)
	}
	logch = make(chan int)
	wait = make(chan bool)
	go Log()
	for i := range num {
		go func(arg int) {
			api := api.API{
				URL: url,
			}

			<-ch[i]
			starttime := Timestamp{}
			starttime.getTime()
			_, res, err := api.GET_ping("ping_test")
			if err != nil {
				fmt.Println(err)
				fmt.Println(res)
			}
			endtime := Timestamp{}
			endtime.getTime()
			logch <- endtime.getms() - starttime.getms()
		}(i)
	}

	time.Sleep(time.Millisecond * 10)

	start := Timestamp{}
	start.getTime()
	for i := range num {
		ch[i] <- true
	}
	<-wait
	end := Timestamp{}
	end.getTime()
	fmt.Println("end time")
	fmt.Println(end.getms() - start.getms())
}
