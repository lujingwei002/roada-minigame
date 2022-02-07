package log

import (
	"log"
	"os"
)

type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

func init() {
	SetLogger(log.New(os.Stderr, "", log.LstdFlags))
	//SetLogger(log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile))
}

var (
	Println func(v ...interface{})
	Printf  func(format string, v ...interface{})
	Fatal   func(v ...interface{})
	Fatalf  func(format string, v ...interface{})
)

func SetLogger(logger Logger) {
	if logger == nil {
		return
	}
	Println = logger.Println
	Printf = logger.Printf
	Fatal = logger.Fatal
	Fatalf = logger.Fatalf
}
