package tools

import (
	"log"
	"os"
	"errors"
	"fmt"
)

type logger struct {
	debug *log.Logger
	warn  *log.Logger
	info  *log.Logger
	err   *log.Logger
}

var Log *logger

func init() {
	var debug, warn, info, errLog *log.Logger

	path := os.Getenv("BC_DEBUG_LOG")
	if d, err := os.Open(path); nil == err {
		debug = log.New(d, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		debug = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	}
	path = os.Getenv("BC_WARN_LOG")
	if w, err := os.Open(path); nil == err {
		warn = log.New(w, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		warn = log.New(os.Stderr, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	path = os.Getenv("BC_INFO_LOG")
	if i, err := os.Open(path); nil == err {
		info = log.New(i, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
	}

	path = os.Getenv("BC_ERR_LOG")
	if e, err := os.Open(path); nil == err {
		errLog = log.New(e, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		errLog = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	Log = &logger{debug, warn, info, errLog}
}

func (l *logger) Debug(v ...interface{}) {
	l.debug.Println(v...)
}

func (l *logger) Warning(v ...interface{}) {
	l.warn.Println(v...)
}

func (l *logger) Info(v ...interface{}) {
	l.info.Println(v...)
}

func (l *logger) Error(v ...interface{}) {
	l.err.Panic(v...)
}

func (l *logger) Get(name string) (*log.Logger, error){
	switch name {
	case "info":
		return l.info, nil
	case "debug":
		return l.debug, nil
	case "warning":
		return l.warn, nil
	case "error":
		return l.err, nil
	default:
		msg := fmt.Sprintf("Logger with name %s doesn't exists", name)
		return nil, errors.New(msg)
	}
}
