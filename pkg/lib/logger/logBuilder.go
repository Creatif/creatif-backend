package logger

import (
	"errors"
	"fmt"
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

	return errors.New(fmt.Sprintf("Could not determine type of message. '%s' given.", t))
}

func NewLogBuilder() LogBuilder {
	return &logBuilder{
		messages:        make(map[string]string),
		equalKeyCounter: make(map[string]int),
	}
}
