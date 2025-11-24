package main

import (
	"strconv"
)

func sendMessage(out Output, senderID int, target, text string) {
	switch target {
	case "0", "server", "host", "console":
		sendToServer(out, senderID, text)
	}
}

func sendToServer(out Output, senderID int, text string) {
	PrintMsg(getName(senderID), "Server", text)
	SendMsg(out, "You", "Server", text)
}

func sendToAll(senderID int, text string) {
	clLock.Lock()
	defer clLock.Unlock()
	for id, c := range clients {
		if id == senderID {
			SendMsg(c, "You", "All", text)
		} else {
			SendMsg(c, getName(senderID), "All", text)
		}
	}
	PrintMsg(getName(senderID), "All", text)
}

func sendToUser(out Output, senderID int, target string, text string) {
	var receiverID int
	if id, err := strconv.Atoi(target); err != nil {
		receiverID = id
	} else {
		tUser := getUser(target)
		if tUser == nil {
			out.WriteLine("Error: user not found")
			return
		}
		receiverID = tUser.ID
	}

	clLock.Lock()
	recv := clients[receiverID]
	clLock.Unlock()

	if recv == nil {
		out.WriteLine("Error: receiver is not connected")
		return
	}

	if receiverID == senderID {
		out.WriteLine("Error: cannot send message to yourself")
		return
	}

	SendMsg(out, "You", getName(receiverID), text)
	SendMsg(recv, getName(senderID), "You", text)
	PrintMsg(getName(senderID), getName(receiverID), text)
}
