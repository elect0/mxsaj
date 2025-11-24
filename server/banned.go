package main

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type BannedStore struct {
	UIDs map[int]bool `json:"uids"`
}

var (
	banned           = BannedStore{UIDs: map[int]bool{}}
	banLock          sync.Mutex
	bannedAttempts   = make(map[int]time.Time)
	bannedAttemptsMu sync.Mutex
	rateLimitDelay   = 5 * time.Second
)

func loadBanned() {
	data, err := os.ReadFile("banned.json")
	if err == nil {
		_ = json.Unmarshal(data, &banned)
	}
}

func saveBanned() {
	data, _ := json.MarshalIndent(banned, "", " ")
	_ = os.WriteFile("banned.json", data, 0644)
}

func isBanned(uid int) bool {
	banLock.Lock()
	defer banLock.Unlock()
	return banned.UIDs[uid]
}

func shouldLogBanned(uid int) bool {
	bannedAttemptsMu.Lock()
	defer bannedAttemptsMu.Unlock()

	now := time.Now()
	last, exists := bannedAttempts[uid]
	if exists && now.Sub(last) < rateLimitDelay {
		return false
	}
	bannedAttempts[uid] = now
	return true
}

func setBanned(uid int, value bool) {
	banLock.Lock()
	defer banLock.Unlock()
	if value {
		banned.UIDs[uid] = true
	} else {
		delete(banned.UIDs, uid)
	}
	saveBanned()
}
