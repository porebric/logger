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

type SenderReader struct {
	errCh  chan errorMsg
	sender Sender
	logger *Logger
}

func NewSender(log *Logger, buffer int, sender Sender) *SenderReader {
	return &SenderReader{
		errCh:  make(chan errorMsg, buffer),
		logger: log,
		sender: sender,
	}
}

func (r *SenderReader) write(ctx context.Context, msg errorMsg) {
	select {
	case r.errCh <- msg:
	default:
		FromContext(ctx).Warn("error queue is full, dropping message")
	}
}

func (r *SenderReader) Run(ctx context.Context) {
	for em := range r.errCh {
		if err := r.sender.Send(ctx, em.err, em.msg, em.kvs...); err != nil {
			r.logger.Error(err, "failed to send error message")
		}
	}
}

func (r *SenderReader) Shutdown(_ context.Context) error {
	return nil
}

func (r *SenderReader) Name() string {
	return "logger sender"
}
