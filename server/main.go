package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"golang.org/x/crypto/ssh"
)

func main() {
	loadUsers()
	loadBanned()

	privateBytes, err := os.ReadFile("server.key")
	if err != nil {
		fmt.Println("Server key not found, generate with:")
		fmt.Println("ssh-keygen -t ed25519 -f server.key")
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		panic(err)
	}

	config := &ssh.ServerConfig{
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			pub := string(ssh.MarshalAuthorizedKey(key))

			userLock.Lock()
			var user *User
			for i := range userStore.Users {
				if userStore.Users[i].Key == pub {
					user = &userStore.Users[i]
					break
				}
			}
			if user == nil {
				user = &User{
					ID:  nextUID(),
					Key: pub,
				}
				userStore.Users = append(userStore.Users, *user)
				saveUsers()
			}
			userLock.Unlock()

			return &ssh.Permissions{
				Extensions: map[string]string{
					"id":   strconv.Itoa(user.ID),
					"name": user.Name,
				},
			}, nil
		},
	}

	config.AddHostKey(private)

	port := "4756"

	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Started Mxsaj on port %s! Server ID is 0\n", port)

	go startConsole()

	for {
		nConn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		sshConn, chans, reqs, err := ssh.NewServerConn(nConn, config)
		if err != nil {
			fmt.Println("SSH handshake failed:", err)
			continue
		}

		userID, _ := strconv.Atoi(sshConn.Permissions.Extensions["id"])
		client := getUser(userID)
		if client == nil {
			fmt.Printf("Error: user %d not found\n", userID)
		}

		if isBanned(userID) {
			fmt.Printf("Banned user %s (%d) tried to connect\n", client.Name, userID)

			for newChannel := range chans {
				newChannel.Reject(ssh.Prohibited, "You are banned from this server")
			}
			sshConn.Close()
			continue
		}

		fmt.Printf("Authenticated %s (%d)\n", getName(userID), userID)
		go ssh.DiscardRequests(reqs)

		go func(uid int, chans <-chan ssh.NewChannel) {
			for newChannel := range chans {
				if newChannel.ChannelType() != "session" {
					newChannel.Reject(ssh.UnknownChannelType, "only session channels are supported")
					continue
				}

				clLock.Lock()
				if _, exists := clients[uid]; exists {
					clLock.Unlock()
					newChannel.Reject(ssh.Prohibited, "you already have an active session")
					continue
				}
				clLock.Unlock()

				ch, _, _ := newChannel.Accept()
				clLock.Lock()
				clients[uid] = ch
				clLock.Unlock()

				go handleClient(ch, uid)
			}
		}(userID, chans)
	}
}
