package main

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	clients = make(map[int]ssh.Channel)
	clLock  sync.Mutex
)

func handleClient(ch ssh.Channel, uid int) {
	out := ChannelOutput{ch: ch}

	defer func() {
		clLock.Lock()
		if c, ok := clients[uid]; ok {
			delete(clients, uid)
			c.Close()
		}
		clLock.Unlock()
		fmt.Printf("%s (%d) disconnected at %s\n", getName(uid), uid, time.Now().Format("15:23"))
		broadcastLeave(uid)
	}()

	userLock.Lock()
	user := getUser(uid)
	userLock.Unlock()

	if user.Name == "" {
		changeName(ch, uid, true)
	}

	fmt.Printf("New connection: %s (%d) at %s\n", user.Name, uid, time.Now().Format("15:23"))
	broadcastJoin(uid)
	viewOnline(ChannelOutput{ch: ch})

	buf := make([]byte, 2048)

	for {
		n, err := ch.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Read error:", err)
			continue
		}

		msg := strings.TrimSpace(string(buf[:n]))
		if msg == "" || isMessageTooLong(msg) {
			if isMessageTooLong(msg) {
				out.WriteLine("Error: message is too long (max 2048 chars)")
			}
			continue
		}

		HandleCommand(out, uid, msg)
	}
}
