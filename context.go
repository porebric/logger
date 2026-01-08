package logger

import (
	"context"

	"github.com/rs/zerolog"
)

type ctxKey string

const senderReaderCtxKey ctxKey = "sender_reader"

func ToContext(ctx context.Context, logger *Logger) context.Context {
	return logger.l.WithContext(ctx)
}

func SenderToContext(ctx context.Context, senderReader *senderReader) context.Context {
	return context.WithValue(ctx, senderReaderCtxKey, senderReader)
}

func FromContext(ctx context.Context) *Logger {
	return &Logger{l: zerolog.Ctx(ctx)}
}

func senderFromContext(ctx context.Context) *senderReader {
	sr, ok := ctx.Value(senderReaderCtxKey).(*senderReader)
	if !ok {
		return nil
	}

	return sr
}
