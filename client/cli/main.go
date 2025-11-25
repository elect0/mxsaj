package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/elect0/mxsaj/client/core"

	"github.com/chzyer/readline"
)

func runCLI (server string, conn io.ReadWriteCloser) {
	var mu sync.Mutex
	rl, err := readline.NewEx(&readline.Config{
		Prompt: "> ",
		InterruptPrompt: "^C",
		EOFPrompt: "exit",
	})

	if err != nil {
		panic(err)
	}

	defer rl.Close()

	fmt.Println("%sConnected to server %s%s\n", core.ColorGreen, server, core.ColorReset)

	go func(){
		buf := make([]byte, 1024)
		for {
			if conn == nil {
				for i := 1; i <= 10; i++ {
					newConn, err := core.ConnectSSH(server)
					if err == nil {
						conn = newConn
						mu.Lock()
				fmt.Printf("%sReconnected!%s\n", core.ColorRed, core.ColorReset)
						mu.Unlock()
						break
					}
					mu.Lock()
					fmt.Printf("%sReconnected attempt %d failed: %v%s\n", core.ColorRed, i, err, core.ColorReset)
					mu.Unlock()
					time.Sleep(2 * time.Second)
				}
				if conn == nil {
					fmt.Printf("%sCannot reconnect after 10 attempts, exiting...%s\n", core.ColorRed, core.ColorReset)
					os.Exit(1)
				}
			}

			n, err := conn.Read(buf)
			if err != nil {
				mu.Lock()
				fmt.Printf("%sServer disconnected, reconnecting...%s\n", core.ColorRed, core.ColorReset)
				mu.Unlock()
				conn.Close()
				conn = nil
				time.Sleep(500 * time.Millisecond)
				continue
			}

			msg := strings.TrimSpace(string[buf[:n]])
			mu.Lock()
			rl.Write([]byte("\r\033[K"))
			rl.Write([]byte(core.GetColor(msg, false) + "\n"))
			rl.Refresh()
			mu.Unlock()
		}
	}()

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		msgToSend := line + "\n"
		
		if _, err = conn.Write([]byte(msgToSend)); err != nil {
			fmt.Printf("%sWrite error: %v%s\n", core.ColorRed, err, core.ColorReset)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:\n  mxsaj generate         # generate SSH keys\n  mxsaj <server:port>    # connect to server")
	}

	cmd := os.Args[1]

	if cmd == "generate" {
		if err := core.GenerateKeys(); err != nil {
			fmt.Println("Error:", err)
		}
		return
	}

	var server string
	if envAddr := os.Getenv("MXSAJ_ADDRESS"); envAddr != "" {
		server =envAddr
	} else {
		server = cmd
	}

	conn, err := core.ConnectSSH(server)
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}

	defer conn.Close()

	runCLI(server, conn)
}
