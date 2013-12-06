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

package main

import (
	"fmt"
	"log"
	"os"
	//"os/exec"
	"bufio"
	"github.com/tobyzxj/goecho/echo"
	"github.com/tobyzxj/goecho/monitor"
	"runtime"
	"sync"
	"time"
)

const (
	SHOW_CLIENT_NUM_MAX int = 30
)

// TCP & UDP echo server with a remoter monitor
func main() {
	var (
		echoclient []*echo.EchoClient
		echoMux    sync.Mutex
	)

	var (
		client_tcp_num int
		client_udp_num int
	)

	// 0. SETUP GOECHO RUNNING ENVIRONMENT
	runtime.GOMAXPROCS(runtime.NumCPU())
	LOGFILE, _ := os.OpenFile("goecho.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	defer LOGFILE.Close()
	log.SetOutput(LOGFILE)
	echo.EchoTimeout = 30 // default is 120 Seconds
	echoclient = make([]*echo.EchoClient, 0)
	go echo.EchoTickCreate(&echoclient, &echoMux) // create a timetick

	// 1. START TCP ECHO
	log.Println("Start TCP echo...")
	go echo.TCPEchoRun(":9999", &echoclient, &echoMux)
	log.Println("Start TCP Succeed.")

	// 2. START UDP ECHO
	log.Println("Start UDP echo...")
	go echo.UDPEchoRun(":9999", &echoclient, &echoMux)
	log.Println("Start UDP Succeed.")

	// 3. START A MONITER
	go monitor.MonitorRun(":60000")

	// 4. Check Cs
	fmt.Printf("\033[2J") // clear screen
	go func() {
		for {
			select {
			case v := <-echo.EchoChan:
				switch v {
				case "CLEAR-SCREEN":
					//fmt.Printf("\033[2J") // clear screen
					fmt.Printf("\033[6;1H")
					for i := 0; i < SHOW_CLIENT_NUM_MAX; i++ {
						fmt.Println("                                                                                           ")
					}
				default:
					log.Println("Unknown Command")
				}

			case <-time.After(10 * time.Second):
				fmt.Printf("\033[2J") // clear screen every 5 seconds
			}
		}
	}()

	go func([]*echo.EchoClient) {
		for {
			fmt.Printf("\033[1;1H")
			fmt.Println("============================== GOECHO SERVER ==============================")
			t := time.Now()
			year, month, day := t.Date()
			hour, min, sec := t.Clock()
			fmt.Printf("goecho - %04d/%02d/%02d %02d:%02d:%02d up, running %d days %02d:%02d:%02d\r\n",
				year, month, day, hour, min, sec,
				(*echo.PEchoTick).Get()/86400,
				(*echo.PEchoTick).Get()/3600,
				(*echo.PEchoTick).Get()/60,
				(*echo.PEchoTick).Get()%60)
			fmt.Printf("total: %d (UDP:%d, TCP:%d) users: %d\r\n",
				client_udp_num+client_tcp_num, client_udp_num, client_tcp_num, monitor.MonitorsGet())
			fmt.Println("===========================================================================")
			fmt.Println("--INDEX----TYPE----REMOTE------------------TIMEOUT-----RECV----------SEND--")
			echoMux.Lock()
			client_tcp_num = 0
			client_udp_num = 0
			for n, v := range echoclient {
				//fmt.Println(n, v)
				if n < SHOW_CLIENT_NUM_MAX {
					fmt.Printf("  %05v    %s     %-21s   %04v        %-10v    %-10v\r\n",
						n, v.EchoType, v.Ip, v.RunTick, v.RecvByteNum, v.SendByteNum)
				}
				if v.EchoType == "TCP" {
					client_tcp_num++
				}
			}
			client_udp_num = len(echoclient) - client_tcp_num
			echoMux.Unlock()
			time.Sleep(200 * time.Millisecond) // delay 200ms
		}
	}(echoclient)

	// command
	// q --> exit goecho
	go func() {
		running := true
		reader := bufio.NewReader(os.Stdin)
		for running {
			data, _, _ := reader.ReadLine()
			cmd := string(data)
			switch cmd {
			case "q":
				log.Println("command", cmd)
				os.Exit(0)
			}
		}
	}()

	// stop here
	select {}
}
