package logger

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

type mockSender struct {
	mu    sync.Mutex
	calls int
	last  errorMsg
	err   error
}

func (m *mockSender) Send(
	_ context.Context,
	err error,
	msg string,
	kvs ...any,
) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.calls++
	m.last = errorMsg{
		err: err,
		msg: msg,
		kvs: kvs,
	}

	return m.err
}

func (m *mockSender) Calls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.calls
}

func (m *mockSender) Last() errorMsg {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.last
}

func Test_SenderReader_Send(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx, buf := getTestData()
	_ = buf // лог нам тут не важен

	ms := &mockSender{}

	sr := NewSender(FromContext(ctx), 1, ms)
	ctx = SenderToContext(ctx, sr)

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go sr.Run(runCtx)

	err := fmt.Errorf("boom")
	Error(ctx, err, "test message", "a", 1)

	// Assert
	require.Eventually(t, func() bool {
		return ms.Calls() == 1
	}, time.Second, 10*time.Millisecond)

	last := ms.Last()
	assert.Equal(t, err, last.err)
	assert.Equal(t, "test message", last.msg)
	assert.Equal(t, []any{"a", 1}, last.kvs)
}

func Test_Error_AfterShutdown_NoPanic(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx, _ := getTestData()

	ms := &mockSender{}

	sr := NewSender(FromContext(ctx), 1, ms)
	ctx = SenderToContext(ctx, sr)

	go sr.Run(ctx)

	_ = sr.Shutdown(ctx)

	// Act + Assert
	require.NotPanics(t, func() {
		Error(ctx, fmt.Errorf("boom"), "after shutdown")
	})
}

func Test_Error_ChannelFull(t *testing.T) {
	t.Parallel()

	// Arrange
	ctx, _ := getTestData()

	ms := &mockSender{}
	sr := NewSender(FromContext(ctx), 1, ms)
	ctx = SenderToContext(ctx, sr)

	// Run НЕ запускаем → канал переполнится
	Error(ctx, fmt.Errorf("first"), "first")

	// Act + Assert
	require.NotPanics(t, func() {
		Error(ctx, fmt.Errorf("second"), "second")
	})
}
