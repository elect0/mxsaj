package main

import (
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

func broadcastJoin(uid int) {
	clLock.Lock()
	defer clLock.Unlock()

	for id, ch := range clients {
		var msg string
		if id == uid {
			msg = fmt.Sprintf("You (%s) joined at %s\nTip: use :help to see available commands\n", getName(uid), time.Now().Format("15:04"))
		} else {
			msg = fmt.Sprintf("%s joined at %s\n", getName(uid), time.Now().Format("15:04"))
		}

		ch.Write([]byte(msg))
	}
}

func broadcastLeave(uid int) {
	clLock.Lock()
	defer clLock.Unlock()

	msg := fmt.Sprintf("%s disconnected\n", getName(uid))
	for _, ch := range clients {
		ch.Write([]byte(msg))
	}

	clearActiveDMByTarget(uid)
}

func broadcastAction(uid int, action string) {
	clLock.Lock()
	clientsCopy := make(map[int]ssh.Channel, len(clients))
	for k, v := range clients {
		clientsCopy[k] = v
	}
	clLock.Unlock()

	for id, ch := range clientsCopy {
		var msg string
		if id == uid {
			msg = fmt.Sprintf("You (%s) %s \n", getName(uid), action)
		} else {
			msg = fmt.Sprintf("%s %s \n", getName(uid), action)
		}

		ch.Write([]byte(msg))
	}
	fmt.Printf("%s (%d) %s\n", getName(uid), uid, action)
}
