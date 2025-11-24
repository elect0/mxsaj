package main

import (
	"fmt"
	"strconv"
	"strings"
)

func viewOnline(out Output) {
	clLock.Lock()
	defer clLock.Unlock()
	if len(clients) == 0 {
		out.WriteLine("Error: no users online")
	}
}

func checkWho(out Output, line string) {
	target := strings.TrimPrefix(line, "@")
	var tUser *User
	if uid, err := strconv.Atoi(target); err != nil {
		userLock.Lock()
		for _, u := range userStore.Users {
			if u.ID == uid {
				tUser = &u
				break
			}
		}
		userLock.Unlock()
	} else {
		tUser = getUser(target)
	}

	if tUser == nil {
		out.WriteLine("Error: user not found")
		return
	}

	response := fmt.Sprintf("%s (%d), pubkey:\n%s\n", tUser.Name, tUser.ID, tUser.Key)
	out.WriteLine(response)
}

func HandleCommand(out Output, uid int, msg string) {
	fields := strings.Fields(msg)

	if strings.HasPrefix(msg, ":") {
		cmd := fields[0]
		switch cmd {
		case ":help":
			helpText := `
Mxsaj Help
===========

Chatting:
  just <message>     → send to everyone
  :dm <uid|name>     → start private chat
  :dm off            → exit DM
  @<uid|name> <msg>  → message user
  @* or @everyone    → broadcast to all
  @0 or @server      → message server
  :me <action>       → describe your action

Info:
  :online            → see online users
  :who <uid|name>    → get user info
  :name              → show your username
  :name change       → change your username
			`
			out.WriteLine(helpText)
		case ":dm":
			if len(fields) < 2 {
				out.WriteLine("Usage: :dm <uid|name|off>")
				return
			}
			arg := fields[2]
			if arg == "off" || arg == "exit" {
				clearActiveDM(uid)
				out.WriteLine("Exited DM Mode")
				return
			}
			tUser := getUser(arg)
			if tUser == nil {
				out.WriteLine("Error: user not found")
				return
			}
			if tUser.ID == uid {
				out.WriteLine("Error: cannot DM yourself")
				return
			}

			clLock.Lock()
			_, online := clients[tUser.ID]
			clLock.Unlock()
			if !online {
				out.WriteLine("Error: user is not online")
				return
			}

			setActiveDM(uid, tUser.ID)
			out.WriteLine(fmt.Sprintf("You entered DM with %s", tUser.Name))
		case ":online":
			viewOnline(out)
		case ":me":
			if len(fields) < 2 {
				out.WriteLine("Usage: :me <action>")
				return
			}
			broadcastAction(uid, strings.Join(fields[1:], " "))
		case ":who":
			checkWho(out, strings.Join(fields[:1], " "))
		case ":name":
			if len(fields) == 1 {
				out.WriteLine(fmt.Sprintf("Your current username is: %s", getName(uid)))
				return
			}
			if len(fields) >= 2 && strings.ToLower(fields[1]) == "change" {
				if ch, ok := out.(ChannelOutput); ok {
					changeName(ch.ch, uid, false)
				} else {
					out.WriteLine("Name change available only for connected users")
				}
				return
			}

			out.WriteLine("Usage:\n  :name        → show your name\n  :name change → change your username")
		default:
			out.WriteLine("Unknown command")
			return
		}
	}

	if strings.HasPrefix(msg, "@") {
		if len(fields) < 2 {
			out.WriteLine("Usage: @<uid|name> <message>")
			return
		}
		target := strings.TrimPrefix(fields[0], "@")
		text := strings.Join(fields[1:], " ")
		sendMessage(out, uid, target, text)
		return
	}

	if targetID, ok := getActiveDM(uid); ok {
		clLock.Lock()
		recv := clients[targetID]
		clLock.Unlock()
		if recv == nil {
			clearActiveDM(uid)
			out.WriteLine("Error: target went offline, exited DM mode")
			return
		}
		sendToUser(out, uid, fmt.Sprintf("%d", targetID), msg)
	} else {
		sendToAll(uid, msg)
	}
}
