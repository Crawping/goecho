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
// under the License.

package monitor

import (
	"bufio"
	"fmt"
	. "github.com/tobyzxj/goecho/echo"
	"sync"
)

var cmdLIST = &Command{
	UsageLine: "LIST [udp/tcp]",
	Short:     "list udp or tcp clients",
}

func init() {
	cmdLIST.Run = listClients
}

func listClients(cmd *Command, args []string, w *bufio.Writer, clients *[]*EchoClient, mux_clients *sync.Mutex) {
	//w.Write([]byte("\r\n"))
	if len(args) == 1 {
		w.Write([]byte("--INDEX----TYPE----REMOTE------------------TIMEOUT-----RECV----------SEND--\r\n"))
		mux_clients.Lock()
		for i, v := range *clients {
			str := fmt.Sprintf("  %05v    %s     %-21s   %04v        %-10v    %-10v\r\n",
				i, v.EchoType, v.Ip, v.RunTick, v.RecvByteNum, v.SendByteNum)
			w.Write([]byte(str))
		}
		mux_clients.Unlock()
	} else if len(args) == 2 && args[1] == "tcp" {
		w.Write([]byte("--INDEX----TYPE----REMOTE------------------TIMEOUT-----RECV----------SEND--\r\n"))
		mux_clients.Lock()
		for i, v := range *clients {
			if v.EchoType == "TCP" {
				str := fmt.Sprintf("  %05v    %s     %-21s   %04v        %-10v    %-10v\r\n",
					i, v.EchoType, v.Ip, v.RunTick, v.RecvByteNum, v.SendByteNum)
				w.Write([]byte(str))
			}
		}
		mux_clients.Unlock()
	} else if len(args) == 2 && args[1] == "udp" {
		w.Write([]byte("--INDEX----TYPE----REMOTE------------------TIMEOUT-----RECV----------SEND--\r\n"))
		mux_clients.Lock()
		for i, v := range *clients {
			if v.EchoType == "UDP" {
				str := fmt.Sprintf("  %05v    %s     %-21s   %04v        %-10v    %-10v\r\n",
					i, v.EchoType, v.Ip, v.RunTick, v.RecvByteNum, v.SendByteNum)
				w.Write([]byte(str))
			}
		}
		mux_clients.Unlock()
	} else {
		w.Write([]byte(cmd.Usage()))
	}

	w.Write([]byte("\r\n"))
	w.Flush()
}
