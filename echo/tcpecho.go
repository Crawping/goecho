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
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// tcp echo server
func TCPEchoRun(service string, stcpc *[]*EchoClient, mux *sync.Mutex) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		log.Println("[TCPERR]", "net.ResolveTCPAddr: ", err.Error())
		return
	}
	tcpl, err := net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		log.Println("[TCPERR]", "net.ListenTCP: ", err.Error())
		return
	}
	defer tcpl.Close()

	for {
		conn, err := tcpl.Accept()
		if err != nil {
			log.Println("[TCPERR]", "tcpl.Accept: ", err.Error())
			return
		}

		// add new udp socket to []*EchoClient
		log.Println("[TCP]", "New TCP client: ->", conn.RemoteAddr())
		mux.Lock()
		tcpc := new(EchoClient)
		tcpc.Id = EchoId
		EchoId++
		tcpc.EchoType = "TCP"
		tcpc.TCPConn = conn
		//conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		tcpc.UDPConn = nil
		tcpc.Ip = fmt.Sprintf("%s", conn.RemoteAddr())
		tcpc.RecvByteNum = 0
		tcpc.SendByteNum = 0
		tcpc.JoinTime = time.Now()
		tcpc.UpdateTime = time.Now()
		tcpc.RunTick = 0
		tcpc.Msg = make(chan string, 10)
		*stcpc = append(*stcpc, tcpc)
		// set tcp timeout
		mux.Unlock()
		EchoChan <- "CLEAR-SCREEN"

		go TCPEchoHandle(conn, stcpc, mux, tcpc)
	}
}

// tcp client handle
func TCPEchoHandle(conn net.Conn, stcpc *[]*EchoClient, mux *sync.Mutex, tcpc *EchoClient) {
	defer func() {
		log.Println("[TCP]", "TCP Client Close: <-", conn.RemoteAddr())
		conn.Close()

		// remove client from []*EchoClient
		mux.Lock()
		for i, v := range *stcpc {
			if v == tcpc {
				if len(*stcpc) == 1 {
					// remove only one
					*stcpc = make([]*EchoClient, 0)
				} else if i == len(*stcpc)-1 {
					// remove last one
					*stcpc = (*stcpc)[:i]
				} else if i == 0 {
					// remove first one
					*stcpc = (*stcpc)[1:]
				} else {
					*stcpc = append((*stcpc)[:i], (*stcpc)[(i+1):]...)
				}
				break
			}
		}
		mux.Unlock()
		EchoChan <- "CLEAR-SCREEN"
	}()

	var buf [2048]byte

	// ECHO
	for {
	LOOP:
		select {
		case v := <-(*tcpc).Msg:
			log.Println("[MSG] TCP CLOSE", v)
			switch v {
			case "CLOSE":
				//log.Println("[MSG] TCP CLOSE")
				goto LOOP
			}

		default:
			// read
			conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(EchoTimeout)))
			n, err := conn.Read(buf[:])
			if err != nil {
				log.Println("[TCPERR]", "conn.Read: ", err.Error())
				if err == io.EOF {
					return
				}

				return
			}

			log.Println("[TCP]", conn.RemoteAddr(), "Recv:", n, "Bytes")
			log.Println("[TCP]", conn.RemoteAddr(), buf[:n])

			// update UpdateTime
			mux.Lock()
			tcpc.RecvByteNum += uint64(n)
			tcpc.UpdateTime = time.Now()
			tcpc.RunTick = 0
			mux.Unlock()
			//EchoChan <- "CLEAR-SCREEN" // not need

			// write
			n, err = conn.Write(buf[:n])
			if err != nil {
				log.Println("[TCPERR]", "conn.Write: ", err.Error())
				return
			}
			log.Println("[TCP]", conn.RemoteAddr(), "Send:", n, "Bytes")
			log.Println("[TCP]", conn.RemoteAddr(), buf[:n])

			mux.Lock()
			tcpc.SendByteNum += uint64(n)
			mux.Unlock()
			//EchoChan <- "CLEAR-SCREEN" // not need
		}
	}
}
