package mylog

import (
	"fmt"
	"sync"
	"time"
)

var mu sync.Mutex
var isDebug bool

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func SetDebug(d bool) {
	isDebug = d
}

func Error(n uint64, s string) {
	defer mu.Unlock()
	mu.Lock()
	if isDebug {
		fmt.Printf("%sError(%d): %s\t%s%s\n", Red, n, time.Now().Local().Format("02.01.2006 15:04:05"), s, Reset)
	}
}

func Info(n uint64, s string) {
	defer mu.Unlock()
	mu.Lock()
	if isDebug {
		fmt.Printf("%sInfo(%d): %s\t%s%s\n", Green, n, time.Now().Local().Format("02.01.2006 15:04:05"), s, Reset)
	}
}
