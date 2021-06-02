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
	file := l.file
	if level == ERROR {
		file = os.Stderr
	}

	if l.level >= level {
		fmt.Fprintf(file, level.String() + format + "\n", args...)
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
	file: os.Stdout,
}
