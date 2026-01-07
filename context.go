package logger

import (
	"context"

	"github.com/rs/zerolog"
)

type ctxKey string

const senderReaderCtxKey ctxKey = "sender_reader"

func ToContext(ctx context.Context, logger *Logger) context.Context {
	return logger.l.WithContext(context.WithValue(ctx, senderReaderCtxKey, logger.senderReader))
}

func FromContext(ctx context.Context) *Logger {
	sr, ok := ctx.Value(senderReaderCtxKey).(*senderReader)
	if !ok {
		sr = nil
	}
	return &Logger{l: zerolog.Ctx(ctx), senderReader: sr}
}
