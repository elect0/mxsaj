package main


import (
	"fmt"
	"io"
	"sync"

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

	fmt.Println("%sConnected to server %s%s\n", rl.Operation.CompleteRefresh)
}
