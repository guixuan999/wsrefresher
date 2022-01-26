package utils

import (
	"fmt"
	"io"
	"os"
	"time"
)

var _loggers = make(map[string]*logger)

type logger struct {
	F io.Writer
}

func (l *logger) Info(format string, params ...interface{}) {
	now := time.Now()
	now_str := now.Format("2006-01-02 15:04:05") // Fucked by go, so we need to use this magic patten "2006-01-02 15:04:05"
	now_str = fmt.Sprintf("%s.%03d", now_str, now.Nanosecond()/(1000*1000))

	format = now_str + " " + format

	fmt.Fprintf(l.F, format, params...)
}

func GetLogger(fn string) *logger {
	if _, ok := _loggers[fn]; !ok {
		f, _ := os.OpenFile(fn, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		_loggers[fn] = &logger{F: f}
	}

	return _loggers[fn]
}
