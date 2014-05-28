package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"strings"

	"code.google.com/p/go.net/websocket"
	"github.com/bobappleyard/readline"
)

var addr = flag.String("s", "localhost:8080", "address of server (host:port)")

func main() {
	log.SetFlags(0)
	flag.Parse()
	if len(strings.SplitN(*addr, ":", 2)) < 2 {
		log.Fatalln("Server address must be of the form host:port")
	}

	origin := "http://" + strings.SplitN(*addr, ":", 2)[0] + "/"
	sock, err := websocket.Dial("ws://"+*addr+"/ws", "", origin)
	if err != nil {
		log.Fatalln("Could not connect to server:", err)
	}

	in := make(chan string)
	out := make(chan string)
	go recvMsgs(sock, in)
	go readCmds(out)

event_loop:
	for {
		select {
		case msg, ok := <-in:
			if !ok {
				readline.Cleanup()
				fmt.Println("\nServer died!")
				break event_loop
			}

			fmt.Print("\033[2K\033[0G")
			fmt.Println("\033[32m" + msg + "\033[39m")
			fmt.Print("> ")

		case cmd, ok := <-out:
			if !ok {
				fmt.Println()
				break event_loop
			}

			if err := websocket.Message.Send(sock, cmd); err != nil {
				log.Println("Could not send command to server:", err)
				continue
			}
		}
	}
}

func readCmds(ch chan string) {
	for {
		cmd, err := readline.String("> ")
		if err == io.EOF {
			close(ch)
			break
		}

		if err != nil {
			log.Fatalln("Error occured while reading command:", err)
		}

		if cmd == "" {
			continue
		}

		readline.AddHistory(cmd)
		ch <- cmd
	}
}

func recvMsgs(sock *websocket.Conn, ch chan string) {
	for {
		var msg string
		err := websocket.Message.Receive(sock, &msg)
		if err == io.EOF {
			close(ch)
			return
		}

		if err != nil {
			ch <- "Could not receive response from server:" + err.Error()
			continue
		}
		ch <- msg
	}
}
