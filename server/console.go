package main

import (
	"strings"

	"github.com/chzyer/readline"
)

func startConsole() {
	go func() {
		rl, err := readline.NewEx(&readline.Config{
			Prompt:      "> ",
			HistoryFile: "console.log",
		})
		if err != nil {
			panic(err)
		}

		defer rl.Close()

		out := ConsoleOutput{}
		out.WriteLine("Mxsaj Server Console started. Type :help for commands.")

		for {
			line, err := rl.Readline()
			if err != nil {
				break
			}
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			HandleCommand(out, 0, line)
		}
	}()
}
