package logger

import (
	"context"
	"github.com/rs/zerolog"
)

type Sender interface {
	Send(ctx context.Context, err error, msg string, kvs ...any) error
}

type senderReader struct {
	ctx    context.Context
	errCh  chan errorMsg
	sender Sender
}

type errorMsg struct {
	err error
	msg string
	kvs []any
}

func (r *senderReader) startErrorSender(l *zerolog.Logger) {
	for em := range r.errCh {
		if err := r.sender.Send(r.ctx, em.err, em.msg, em.kvs...); err != nil {
			l.Error().Err(err).Msg("failed to send error message")
		}
	}
}
