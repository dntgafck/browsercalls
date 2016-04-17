package tools

import (
	"log"
	"os"
)

type Logger struct {
	debug *log.Logger
	warn  *log.Logger
	info  *log.Logger
	err   *log.Logger
}

var Log *Logger

func init() {
	var debug, warn, info, err *log.Logger

	path := os.Getenv("BC_DEBUG_LOG")
	if d, err := os.Open(path); nil == err {
		debug = log.New(d, "DEBUG", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		debug = log.New(os.Stdout, "DEBUG", log.Ldate|log.Ltime|log.Lshortfile)
	}
	path = os.Getenv("BC_WARN_LOG")
	if w, err := os.Open(path); nil == err {
		debug = log.New(w, "WARNING", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		debug = log.New(os.Stderr, "WARNING", log.Ldate|log.Ltime|log.Lshortfile)
	}

	path = os.Getenv("BC_INFO_LOG")
	if i, err := os.Open(path); nil == err {
		debug = log.New(i, "INFO", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		debug = log.New(os.Stdout, "INFO", log.Ldate|log.Ltime|log.Lshortfile)
	}

	path = os.Getenv("BC_ERR_LOG")
	if e, err := os.Open(path); nil == err {
		debug = log.New(e, "ERROR", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		debug = log.New(os.Stderr, "ERROR", log.Ldate|log.Ltime|log.Lshortfile)
	}

	Log = &Logger{debug, warn, info, err}
}

func (l *Logger) Debug(v ...interface{}) {
	l.debug.Println(v...)
}

func (l *Logger) Warining(v ...interface{}) {
	l.warn.Println(v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.info.Println(v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.err.Panic(v...)
}
