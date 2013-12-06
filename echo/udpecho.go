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
	"log"
	"net"
	"sync"
	"time"
)

// udp echo server
func UDPEchoRun(service string, sudpc *[]*EchoClient, mux *sync.Mutex) {
	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	if err != nil {
		log.Println("[UDPERR]", "net.ResolveUDPAddr: ", err.Error())
		return
	}
	udpconn, err := net.ListenUDP("udp4", udpAddr)
	if err != nil {
		log.Println("[UDPERR]", "net.ListenUDP: ", err.Error())
		return
	}
	defer udpconn.Close()

	UDPEchoHandle(&udpconn, sudpc, mux)
}

// udp client handle
func UDPEchoHandle(conn **net.UDPConn, sudpc *[]*EchoClient, mux *sync.Mutex) {

	var buf [2048]byte
	var isFind bool = false
	var udp_ip string

	// udp echo handle
	for {
		n, remote_addr, err := (*conn).ReadFromUDP(buf[0:])
		if err != nil {
			log.Println("[UDPERR]", "udp receiveFrom ", remote_addr, "faild: ", err.Error())
			continue
		}
		log.Println("[UDP]", remote_addr, "Recv:", n, "Bytes")
		log.Println("[UDP]", remote_addr, buf[:n])

		// Check whether remote_addr already exists
		isFind = false
		udp_ip = fmt.Sprintf("%s:%d", remote_addr.IP, remote_addr.Port) // format
		mux.Lock()
		for i, v := range *sudpc {
			if v.Ip == udp_ip {
				// remote_addr already exists
				isFind = true
				// update RecvByteNum & UpdateTime
				(*sudpc)[i].RecvByteNum += uint64(n)
				(*sudpc)[i].UpdateTime = time.Now()
				(*sudpc)[i].RunTick = 0
				//EchoChan <- "CLEAR-SCREEN" // not need
				break
			}
		}
		mux.Unlock()

		// add new udp socket to []*EchoClient
		if !isFind {
			log.Println("[UDP]", "New UDP client: ->", remote_addr)
			mux.Lock()
			udpc := new(EchoClient)
			udpc.Id = EchoId
			EchoId++
			udpc.EchoType = "UDP"
			udpc.TCPConn = nil
			udpc.UDPConn = *conn
			udpc.Ip = fmt.Sprintf("%s:%v", remote_addr.IP, remote_addr.Port)
			udpc.RecvByteNum = uint64(n)
			udpc.SendByteNum = 0
			udpc.JoinTime = time.Now()
			udpc.UpdateTime = time.Now()
			udpc.RunTick = 0
			*sudpc = append(*sudpc, udpc)
			mux.Unlock()
			EchoChan <- "CLEAR-SCREEN"
		}

		// udp echo
		n, err = (*conn).WriteTo(buf[:n], remote_addr)
		if err != nil {
			log.Println("[UDPERR]", "udp sendto ", remote_addr, "faild: ", err.Error())
			continue
		}
		log.Println("[UDP]", remote_addr, "Send:", n, "Bytes")
		log.Println("[UDP]", remote_addr, buf[:n])

		// update RecvByteNum
		mux.Lock()
		for i, v := range *sudpc {
			if v.Ip == udp_ip {
				// found
				(*sudpc)[i].SendByteNum += uint64(n)
				//EchoChan <- "CLEAR-SCREEN" // not need
				break
			}
		}
		mux.Unlock()
	}
}
