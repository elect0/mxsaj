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

	g
}
