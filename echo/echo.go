/**
 * Created with IntelliJ IDEA.
 * User: toby.zxj
 * Email: toby.zxj@gamil.com
 * Date: 13-11-24 下午2:44
 */
package echo

import (
	//"fmt"
	"net"
	//"runtime"
	"sync"
	"time"
)

var EchoId int = 0

//var EchoClearScreen bool = true
var EchoTimeout uint64 = 120
var EchoChan = make(chan string, 10)
var PEchoTick *EchoTick = nil

// 120 Seconds

// echo client, tcp & udp client
type EchoClient struct {
	mux         sync.RWMutex
	Id          int
	EchoType    string
	TCPConn     net.Conn
	UDPConn     *net.UDPConn
	Ip          string
	RecvByteNum uint64
	SendByteNum uint64
	JoinTime    time.Time `orm:"index"`
	UpdateTime  time.Time `orm:"index"`
	RunTick     uint64
	Msg         chan string
}

type EchoTick struct {
	tick uint64
}

func (this *EchoTick) Get() uint64 {
	return this.tick
}

func (this *EchoTick) Tick() {
	this.tick++
}

func (this *EchoTick) Set(tick uint64) {
	this.tick = tick
}

func EchoTickCreate(client *[]*EchoClient, mux *sync.Mutex) {
	PEchoTick = &EchoTick{tick: 0}
	timer := time.NewTicker(1*time.Second)

	for {
		select {
		case <-timer.C:
			EchoTickTick(client, mux, PEchoTick)
		}
	}
}

// performed once every second
func EchoTickTick(client *[]*EchoClient, mux *sync.Mutex, tick *EchoTick) {

	if tick != nil {
		tick.Tick()
	}

	var index int = 0

	mux.Lock()
RECHECK:
	for i, v := range *client {
		if i >= index {
			((*client)[i]).RunTick++
		} else {
			continue
		}

		// Check timeout, remove it
		if v.RunTick > EchoTimeout {
			index = i
			// panic: runtime error: index out of range
			switch v.EchoType {
			case "TCP":
				// Do nothing
				//v.Msg <- "CLOSE"
				//v.TCPConn.Close()
				//fallthrough
				//runtime.Goexit()

			case "UDP":
				if len(*client) == 1 {
					// remove only one
					*client = make([]*EchoClient, 0)
				} else if i == len(*client) - 1 {
					// remove last one
					*client = (*client)[:i]
					break RECHECK
				} else if i == 0 {
					// remove first one
					*client = (*client)[1:]
				} else {
					*client = append((*client)[:i], (*client)[(i + 1):]...)
				}
				//EchoChan <- "CLEAR-SCREEN"
			default:
			}
			goto RECHECK
		}
	}
	mux.Unlock()
}
