package bplog

import (
	"log"
	"os"
)

var bplog *log.Logger

func init() {
	bplog = log.New(os.Stdout, "[buildpack-aliyun Action] ", 0)
}

func Println(v ...interface{}) {
	bplog.Println(v)
}

func Printf(format string, v ...interface{}) {
	bplog.Printf(format, v...)
}

func Fatal(v ...interface{}) {
	bplog.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	bplog.Fatalf(format, v...)
}
