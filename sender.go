package logger

import (
	"context"
)

type Sender interface {
	Send(ctx context.Context, err error, msg string, kvs ...any) error
}

type errorMsg struct {
	err error
	msg string
	kvs []any
}

type senderReader struct {
	errCh  chan errorMsg
	sender Sender
	logger *Logger
}

func NewSender(log *Logger, buffer int, sender Sender) *senderReader {
	return &senderReader{
		errCh:  make(chan errorMsg, buffer),
		logger: log,
		sender: sender,
	}
}

func (r *senderReader) write(ctx context.Context, msg errorMsg) {
	select {
	case r.errCh <- msg:
	default:
		FromContext(ctx).Warn("error queue is full, dropping message")
	}
}

func (r *senderReader) Run(ctx context.Context) {
	for em := range r.errCh {
		if err := r.sender.Send(ctx, em.err, em.msg, em.kvs...); err != nil {
			r.logger.Error(err, "failed to send error message")
		}
	}
}

func (r *senderReader) Shutdown(_ context.Context) error {
	return nil
}
