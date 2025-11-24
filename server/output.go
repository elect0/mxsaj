package main

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

type Output interface {
	WriteLine(msg string)
}

type ChannelOutput struct {
	ch ssh.Channel
}

func (c ChannelOutput) WriteLine (msg string) {
	c.ch.Write([]byte(msg + "\n"))
}

type ConsoleOutput struct {}

func (ConsoleOutput) WriteLine (msg string) {
	fmt.Println(msg)
}
