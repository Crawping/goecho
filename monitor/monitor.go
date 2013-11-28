/**
 * Created with IntelliJ IDEA.
 * User: toby.zxj
 * Email: toby.zxj@gamil.com
 * Date: 13-11-24 下午8:34
 */
package monitor

import (
	"fmt"
	"net"
	"log"
	"sync"
	"container/list"
	"time"
	"io"
	"bufio"
//	"net/textproto"
)

var MonitorNum int = 0
var monitorList *list.List

const (
	MONITOR_WELLCOME string = "=========================\r\n" +
	                          "Wellcome to Goecho Server\r\n" +
	                          "=========================\r\n"
	MONITOR_USERNAME string = "Name: "
	MONITOR_PASSWD string = "Passwd: "

)

// Element of monitor client
type MonitorClient struct {
	mux         sync.RWMutex
	Id          int
	conn        net.Conn
	Ip          string
	RunTick     uint64
	Msg         chan string
}

func MonitorRun(service string) {
	var (
		monitorId int = 0
		monitorMux sync.Mutex
	)
	// create a list to save monitorClient
	monitorList = list.New()

	// create a TCP server listen
	addr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		log.Println("[MONITORERR]", "net.ResolveTCPAddr: ", err.Error())
		return
	}
	tcpl, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		log.Println("[MONITORERR]", "net.ListenTCP", err.Error())
		return
	}
	defer tcpl.Close()

	go func(l *list.List, mux *sync.Mutex) {
		for {
			var i int = 0
			mux.Lock()
			// iterate
			for e := l.Front(); e != nil; e = e.Next() {
				// do something with e.Value
				i++
			}
			MonitorNum = l.Len()
			mux.Unlock()

			time.Sleep(5*time.Second)
		}
	}(monitorList, &monitorMux)


	for {
		conn, err := tcpl.Accept()
		if err != nil {
			log.Println("[MONITORERR]", "net.Accept: ", err.Error())
			return
		}

		// add a monitorClient
		monitorMux.Lock()
		elem := new(MonitorClient)
		elem.Id = monitorId
		monitorId++
		elem.conn  = conn
		elem.Ip = fmt.Sprintf("%s", conn.RemoteAddr())
		elem.RunTick = 0
		elem.Msg = make(chan string, 10)
		monitorList.PushBack(elem)
		monitorMux.Unlock()

		go monitorHandle(monitorList, &monitorMux, elem)
	}
}

// monitor client handle
func monitorHandle(l *list.List, mux *sync.Mutex, c *MonitorClient) {
	defer func() {
		c.conn.Close()

		// remove e
		mux.Lock()
		for e := l.Front(); e != nil; e = e.Next() {
			// do something with e.Value
			if e.Value == c  {
				l.Remove(e)
			}
		}
		mux.Unlock()
	}()

	stream_in := bufio.NewReaderSize(c.conn, 64) // max command size < 64
	stream_out := bufio.NewWriterSize(c.conn, 2048)

	// say hello to user
	_, err := stream_out.Write([]byte(MONITOR_WELLCOME))
	_, err = stream_out.Write([]byte(MONITOR_USERNAME))
	err = stream_out.Flush()
	if err != nil {
		return
	}


	// read username & passwd
	c.conn.SetReadDeadline(time.Now().Add(time.Second*time.Duration(120))) // 120s timeout
	data, _, err := stream_in.ReadLine()
	CheckMonitorError(err)
	if err != nil {
		return
	}

	_, err = stream_out.Write([]byte(MONITOR_PASSWD))
	err = stream_out.Flush()
	if err != nil {
		return
	}

	username := string(data)
	c.conn.SetReadDeadline(time.Now().Add(time.Second*time.Duration(120))) // 120s timeout
	data, _, err = stream_in.ReadLine()
	CheckMonitorError(err)
	if err != nil {
		return
	}
	passwd := string(data)

	// check username & passwd
	if username != "admin" ||  passwd != "admin" {
		return
	}

	for {
		_, err = stream_out.Write([]byte("\r\ngoecho>> "))
		err = stream_out.Flush()
		if err != nil {
			return
		}

		// read
		c.conn.SetReadDeadline(time.Now().Add(time.Second*time.Duration(120))) // 120s timeout
		data, _, err = stream_in.ReadLine()
		CheckMonitorError(err)
		if err != nil {
			return
		}
		cmd := string(data)
		switch cmd {
		case "list":
			_, err = stream_out.Write([]byte("HELLO\r\n"))
		case "exit":
			_, err = stream_out.Write([]byte("ByeBye\r\n"))
			err = stream_out.Flush()
			return
		default :
			_, err = stream_out.Write([]byte("Unknown Command\r\n"))
		}
		err = stream_out.Flush()
		if err != nil {
			return
		}
	}
}

func MonitorsGet() int {
	if monitorList == nil {
		return 0
	} else {
		return monitorList.Len()
	}

	panic("unreachable");
}

func CheckMonitorError(err error) error {
	if err != nil {
		log.Println("[MONITORERR]", "stream_in.ReadLine: ", err.Error())

		if err == io.EOF {
			return err
		}

		return err
	}

	return nil
}

//////////////// MAP Command --> Handle ///////////////////////

