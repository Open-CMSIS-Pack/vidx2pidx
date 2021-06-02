package main

import (
	"io"
	"fmt"
	"os"
)

type LevelType int

const (
	ERROR LevelType = iota
	INFO
	DEBUG
)

func (l LevelType) String() string {
	return [...]string{"E: ", "I: ", "D: "}[l]
}

type LoggerType struct {
	level LevelType
	file io.Writer
}

func (l *LoggerType) output(level LevelType, format string, args ...interface{}) {
	if l.level >= level {
		fmt.Fprintf(l.file, level.String() + format + "\n", args...)
	}
}

func (l *LoggerType) Debug(format string, args ...interface{}) {
	l.output(DEBUG, format, args...)
}

func (l *LoggerType) Info(format string, args ...interface{}) {
	l.output(INFO, format, args...)
}

func (l *LoggerType) Error(format string, args ...interface{}) {
	l.output(ERROR, format, args...)
}

func (l *LoggerType) SetLevel(level LevelType) {
	l.level = level
}

func (l *LoggerType) SetFile(file io.Writer) {
	l.file = file
}

var Logger = LoggerType {
	level: ERROR,
}

func init() {
	if os.Getenv("TESTING") == "1" {
		f, err := os.OpenFile("testing.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		ExitOnError(err)
		Logger.SetFile(f)
	} else {
		Logger.SetFile(os.Stdout)
	}
}
