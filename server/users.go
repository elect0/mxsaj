package main

import (
	"encoding/json"
	"os"
	"sync"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

type UserStore struct {
	Users []User `json:"users"`
}

var (
	userStore = UserStore{Users: []User{}}
	userLock  sync.Mutex
)

func loadUsers() {
	data, err := os.ReadFile("users.json")
	if err == nil {
		_ = json.Unmarshal(data, &userStore)
	}
}

func saveUsers() {
	data, _ := json.MarshalIndent(userStore, "", " ")
	_ = os.WriteFile("users.json", data, 0644)
}

func nextUID() int {
	max := 0
	for _, u := range userStore.Users {
		if u.ID > max {
			max = u.ID
		}
	}
	return max + 1
}
