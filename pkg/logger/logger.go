package logger

import (
	"fmt"
	"strings"
	"time"
)

type Logger interface {
	Info(msg string)
	Warn(msg string, err error)
	Error(msg string, err error)
}

type Log struct {
}

func New() *Log {
	return &Log{}
}

func (l *Log) Info(msg string) {
	fmt.Println("INFO", formatTime(time.Now()), "MESSAGE: "+msg)
}

func (l *Log) Warn(msg string, err error) {
	fmt.Println("WARN", formatTime(time.Now()), "MESSAGE: "+msg, "WARN: "+err.Error())
}

func (l *Log) Error(msg string, err error) {
	fmt.Println("ERROR", formatTime(time.Now()), "MESSAGE: "+msg, "ERROR: "+err.Error())
}

func formatTime(now time.Time) string {
	time := strings.Replace(now.UTC().Format(time.RFC3339), "T", " ", -1)
	return "[" + strings.Replace(time, "Z", "", -1) + "]"
}
