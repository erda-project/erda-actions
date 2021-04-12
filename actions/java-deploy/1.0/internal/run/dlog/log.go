package dlog

import (
	"log"
	"os"
)

var l *log.Logger

func init() {
	l = log.New(os.Stdout, "[Java-Deploy Action] ", 0)
}

func Println(v ...interface{}) {
	l.Println(v)
}

func Printf(format string, v ...interface{}) {
	l.Printf(format, v...)
}

func Fatal(v ...interface{}) {
	l.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	l.Fatalf(format, v...)
}
