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
	"strings"
	"sync"
)

type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	// The w is the output stream.
	// The clients is echo clients.
	// The mux_clients is mux for echo clients.
	Run func(cmd *Command, args []string, w *bufio.Writer, clients *[]*EchoClient, mux_clients *sync.Mutex)

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'go help' output.
	Short string
}

// Name returns the command's name: the first word in the usage line.
func (this *Command) Name() string {
	name := this.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command such as importpath.
func (this *Command) Runnable() bool {
	return this.Run != nil
}

// Usage returns the command's usage infomation
func (this *Command) Usage() string {
	str := "Usage: \r\n\t" + this.UsageLine
	return str
}

// Help returns the command's help infomation
func (this *Command) help() string {
	return this.Short
}
