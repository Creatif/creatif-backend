package logger

import (
	"fmt"
	"os"
)

const info = "info"
const errorLog = "error"
const warn = "warn"

type LogBuilder interface {
	Add(string, string)
	Flush(t string) error
}

type logBuilder struct {
	messages        map[string]string
	equalKeyCounter map[string]int
}

func (l *logBuilder) Add(key, message string) {
	_, ok := l.equalKeyCounter[key]
	if !ok {
		l.equalKeyCounter[key] = 0
	} else {
		l.equalKeyCounter[key] = l.equalKeyCounter[key] + 1
	}

	l.messages[fmt.Sprintf("%s_%d", key, l.equalKeyCounter[key])] = message
}

func (l *logBuilder) Flush(t string) error {
	if os.Getenv("APP_ENV") != "prod" {
		//fmt.Println(l.messages)
	}

	clear := func() {
		l.messages = nil
		l.equalKeyCounter = nil
	}
	defer clear()

	if t == info {
		Info(l.messages)

		return nil
	}

	if t == errorLog {
		Error(l.messages)

		return nil
	}

	if t == warn {
		Warn(l.messages)

		return nil
	}

	Error(l.messages)
	l.messages = nil
	l.equalKeyCounter = nil

	return nil
}

func NewLogBuilder() LogBuilder {
	return &logBuilder{
		messages:        make(map[string]string),
		equalKeyCounter: make(map[string]int),
	}
}
