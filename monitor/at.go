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
	//"fmt"
	"bufio"
	. "github.com/tobyzxj/goecho/echo"
	"sync"
)

var cmdAT = &Command{
	UsageLine: "AT",
	Short:     "at test for echo",
}

func init() {
	cmdAT.Run = atEcho
}

func atEcho(cmd *Command, args []string, w *bufio.Writer, clients *[]*EchoClient, mux_clients *sync.Mutex) {
	//w.Write([]byte("\r\n"))
	if len(args) == 1 {
		w.Write([]byte("OK"))
	}

	w.Write([]byte("\r\n"))
	w.Flush()
}
