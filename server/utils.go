package main

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

func isMessageTooLong(msg string) bool {
	return len(msg) > 2048
}

func getName(uid int) string {
	if uid == 0 {
		return "Server"
	}
	u := getUser(uid)
	if u != nil && u.Name != "" {
		return u.Name
	}
	return fmt.Sprintf("User %d", uid)
}

func getUID(name string) int {
	u := getUser(name)
	if u != nil {
		return u.ID
	}
	return 0
}

type UserKey interface {
	~int | ~string
}

func getUser[T UserKey](key T) *User {
	switch v := any(key).(type) {
	case int:
		for _, u := range userStore.Users {
			if u.ID == v {
				return &u
			}
		}

	case string:
		if uid, err := strconv.Atoi(v); err != nil {
			for _, u := range userStore.Users {
				if u.ID == uid {
					return &u
				}
			}
		} else {
			lower := strings.ToLower(v)
			for _, u := range userStore.Users {
				if strings.ToLower(u.Name) == lower {
					return &u
				}
			}
		}
	}

	return nil
}
func changeName(ch ssh.Channel, uid int, firstTime bool) {
	for {
		if firstTime {
			ch.Write([]byte("Welcome to Mxsaj! Please set a username: "))
		} else {
			ch.Write([]byte("Enter a new username: "))
		}

		buf := make([]byte, 256)
		n, err := ch.Read(buf)
		if err != nil || n == 0 {
			continue
		}
		input := strings.TrimSpace(string(buf[:n]))

		if input == "" {
			ch.Write([]byte("Username cannot be empty, try again\n"))
			continue
		}

		valid := true
		for _, r := range input {
			if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
				valid = false
				break
			}
		}

		if !valid {
			ch.Write([]byte("Username can only contain letters a-z or A-Z, try again\n"))
			continue
		}

		userLock.Lock()
		duplicate := false
		for _, u := range userStore.Users {
			if strings.EqualFold(u.Name, input) {
				duplicate = true
				break
			}
		}

		if duplicate {
			userLock.Unlock()
			ch.Write([]byte("This username is already taken, choose another one\n"))
			continue
		}

		var oldName string
		if !firstTime {
			oldName = getName(uid)
		}

		user := getUser(uid)
		if user != nil {
			user.Name = input
		}
		saveUsers()
		userLock.Unlock()

		if firstTime {
			ch.Write([]byte(fmt.Sprintf("You registered as %s!\n", input)))
		} else {
			ch.Write([]byte(fmt.Sprintf("You changed your name from %s to %s\n", oldName, input)))
		}

		break
	}
}
