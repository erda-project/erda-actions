package log

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/erda-project/erda/pkg/color"
	"github.com/erda-project/erda/pkg/sqlparser/migrator"
)

var (
	infoL, errorL, warnL, debugL log.Logger
)

func init() {
	infoL.SetOutput(os.Stdout)
	errorL.SetOutput(os.Stderr)
	warnL.SetOutput(io.MultiWriter(os.Stdout, os.Stderr))
	debugL.SetOutput(os.Stdout)

	infoL.SetPrefix("[ACTION INFO]")
	errorL.SetPrefix("[ACTION ERROR]")
	warnL.SetPrefix("[ACTION WARN]")
	debugL.SetPrefix("[ACTION DEBUG]")

	migrator.Infof = Infof
	migrator.Infoln = Infoln
	migrator.Errorf = Errorf
	migrator.Errorln = Errorln
	migrator.Warnf = Warnf
	migrator.Warnln = Warnln
	migrator.Debugf = Debugf
	migrator.Debugln = Debugln
	migrator.Fatalf = Fatalf
	migrator.Fatalln = Fatalln
}

func Infof(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	infoL.Println(color.Green(s))
}

func Infoln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	infoL.Println(color.Green(s))
}

func Errorf(format string, v ...interface{}) {
	errorL.Printf(format, v...)
}

func Errorln(v ...interface{}) {
	errorL.Println(v...)
}

func Debugf(format string, v ...interface{}) {
	warnL.Printf(format, v...)
}

func Debugln(v ...interface{}) {
	warnL.Println(v...)
}

func Fatalf(format string, v ...interface{}) {
	errorL.Fatalf(format, v...)
}

func Fatalln(v ...interface{}) {
	errorL.Fatalln(v...)
}

func Warnf(format string, v ...interface{}) {
	warnL.Printf(format, v...)
}

func Warnln(v ...interface{}) {
	warnL.Println(v...)
}
