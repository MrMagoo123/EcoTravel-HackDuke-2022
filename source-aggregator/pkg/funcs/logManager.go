package funcs

import (
	"fmt"
	"github.com/gookit/color"
	"time"
)

func InitiateLogger() {
	go logger()
}

var (
	handlerMap = map[string]func(string, ...interface{}){
		"r": color.Red.Printf,
		"b": color.Blue.Printf,
		"c": color.Cyan.Printf,
		"m": color.Magenta.Printf,
		"g": color.Green.Printf,
	}
)

func logger() {
	for l := range LogChan {
		printLog(l)
	}
}

func printLog(l Log) {
	handlerFunc := handlerMap[l.Color]
	if handlerFunc == nil {
		handlerFunc = color.White.Printf
	}

	color.White.Printf("|")
	color.Gray.Printf(" [%s] ", time.Now().Format("15:04:05.000"))
	color.White.Printf("|")
	if l.ID != "" {
		handlerFunc(" " + l.ID)
		color.White.Printf(" | ")
		handlerFunc(l.Msg)
	} else {
		handlerFunc(" " + l.ID)
		handlerFunc(l.Msg)
	}

	fmt.Println("")
}
