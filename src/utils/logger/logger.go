package logger

import (
	"fmt"
	"go.uber.org/dig"
	"log"
)

type Logger struct {
	ILoggerTransport
	transports []ILoggerTransport
	level      int // 0 - debug, 1 - log, 2 - error
	name       string
}

type ILoggerTransport interface {
	Log(str string, ctx string)
	Debug(str string, ctx string)
	Error(str string, ctx string)
}

func CreateLogger(name string, level int) *Logger {
	return &Logger{transports: []ILoggerTransport{&ConsoleTransport{}}, name: name, level: level}
}

func (r *Logger) Clone(name string) *Logger {
	return CreateLogger(name, r.level)
}

func (r *Logger) Provide(container *dig.Scope) *Logger {
	if err := container.Provide(func() *Logger {
		return r
	}); err != nil {
		log.Fatal(err)
	}

	return r
}

func (r *Logger) Log(str string) {
	if r.level <= 1 {
		for _, v := range r.transports {
			v.Log(str, r.name)
		}
	}
}
func (r *Logger) Debug(str string) {
	if r.level == 0 {
		for _, v := range r.transports {
			v.Debug(str, r.name)
		}
	}
}
func (r *Logger) Error(str string) {
	if r.level <= 3 {
		for _, v := range r.transports {
			v.Error(str, r.name)
		}
	}
}

//

type ConsoleTransport struct {
	ILoggerTransport
}

func (r *ConsoleTransport) Log(str string, ctx string) {
	fmt.Println("["+ctx+"] log -", str)
}
func (r *ConsoleTransport) Debug(str string, ctx string) {
	fmt.Println("["+ctx+"] debug -", str)
}
func (r *ConsoleTransport) Error(str string, ctx string) {
	fmt.Println("["+ctx+"] error -", str)
}
