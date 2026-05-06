package logger

import (
	"log"
	"os"
)

type Logger struct {
	service string
}

func New(service string) *Logger {
	return &Logger{
		service: service,
	}
}

func (l *Logger) Info(msg string) {
	log.Printf(`{"level":"info","service":"%s","msg":"%s"}`, l.service, msg)
}

func (l *Logger) Error(msg string) {
	log.Printf(`{"level":"error","service":"%s","msg":"%s"}`, l.service, msg)
}

func (l *Logger) Fatal(msg string) {
	log.Printf(`{"level":"fatal","service":"%s","msg":"%s"}`, l.service, msg)
	os.Exit(1)
}
