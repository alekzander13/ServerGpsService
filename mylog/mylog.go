package mylog

import (
	"fmt"
	"sync"
	"time"
)

var mu sync.Mutex
var IsDebug bool

func Error(n uint64, s string) {
	defer mu.Unlock()
	mu.Lock()
	if IsDebug {
		fmt.Printf("Error(%d): %s\t%s\n", n, s, time.Now().Local().GoString())
	}
}

func Info(n uint64, s string) {
	defer mu.Unlock()
	mu.Lock()
	if IsDebug {
		fmt.Printf("Info(%d): %s\t%s\n", n, s, time.Now().Local().GoString())
	}
}
