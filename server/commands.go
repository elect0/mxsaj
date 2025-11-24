package main

import (
	"fmt"
	"strconv"
	"strings"
)

func viewOnline (out Output) {
	clLock.Lock()
	defer clLock.Unlock()
	if len(clients) == 0 {
		out.WriteLine("Error: no users online")
	} 

}
