package logger

import "context"

type Sender interface {
	Send(ctx context.Context, err error, msg string, kvs ...any) error
}

type errorMsg struct {
	err error
	msg string
	kvs []any
}

func (l *Logger) startErrorSender() {
	for em := range l.errCh {
		l.sem <- struct{}{}
		go func(em errorMsg) {
			defer func() { <-l.sem }()

			if l.sender != nil {
				if err := l.sender.Send(l.senderCtx, em.err, em.msg, em.kvs...); err != nil {
					l.l.Error().Err(err).Msg("failed to send error message")
				}
			}
		}(em)
	}
}
