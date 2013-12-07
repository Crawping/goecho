// Copyright 2013 toby.zxj
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.package main

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

// Get tick's tick
func (this *EchoTick) Get() uint64 {
	return this.tick
}

// Tick tick
func (this *EchoTick) Tick() {
	this.tick++
}

// Set the value of tick
func (this *EchoTick) Set(tick uint64) {
	this.tick = tick
}

// Create a echo tick
func EchoTickCreate(clients *[]*EchoClient, mux *sync.Mutex) {
	PEchoTick = &EchoTick{tick: 0}
	timer := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-timer.C:
			EchoTickTick(clients, mux, PEchoTick)
		}
	}
}

// performed once every second
func EchoTickTick(clients *[]*EchoClient, mux *sync.Mutex, tick *EchoTick) {

	if tick != nil {
		tick.Tick()
	}

	var index int = 0

	mux.Lock()
RECHECK:
	for i, v := range *clients {
		if i >= index {
			((*clients)[i]).RunTick++
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
				if len(*clients) == 1 {
					// remove only one
					*clients = make([]*EchoClient, 0)
				} else if i == len(*clients)-1 {
					// remove last one
					*clients = (*clients)[:i]
					break RECHECK
				} else if i == 0 {
					// remove first one
					*clients = (*clients)[1:]
				} else {
					*clients = append((*clients)[:i], (*clients)[(i+1):]...)
				}
				//EchoChan <- "CLEAR-SCREEN"
			default:
			}
			goto RECHECK
		}
	}
	mux.Unlock()
}
