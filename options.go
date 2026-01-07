package logger

import (
	"context"
	"github.com/rs/zerolog"
	"os"
)

type Option func(*Logger)

func WithPlainText() Option {
	return func(l *Logger) {
		zl := l.l.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Logger()
		l.l = &zl
	}
}

func WriteToFile(f *os.File) Option {
	return func(l *Logger) {
		zl := l.l.Output(zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout}, f)).With().Logger()
		l.l = &zl
	}
}

func WithErrorSender(senderCtx context.Context, s Sender) Option {
	return func(l *Logger) {
		l.sender = s
		l.senderCtx = senderCtx
	}
}

func test(l *Logger) {
	zl := l.l.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Logger()
	l.l = &zl
}
