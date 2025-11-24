package main

import (
	"sync"
)

var (
	activeDM = make(map[int]int)
	dmLock sync.Mutex
)

func setActiveDM(user, target int) {
	dmLock.Lock()
	activeDM[user] = target
	dmLock.Unlock()
}

func clearActiveDM(user int) {
	dmLock.Lock()
	delete(activeDM, user)
	dmLock.Unlock()
}

func getActiveDM(user int) (int, bool) {
	dmLock.Lock()
	t, ok := activeDM[user]
	dmLock.Unlock()
	return t, ok
}

func clearActiveDMByTarget (target int) {
	dmLock.Lock()
	for u, t := range activeDM {
		if t == target {
			delete(activeDM, u)
		}
	}
	dmLock.Unlock()
}
