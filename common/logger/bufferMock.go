package logger

import (
	"bytes"
	"strings"
	"sync"
	"testing"
	"time"
)

// BufferMock is a buffer on which you can wait for a certain log
// to be written. Inject it to your logger in your tests.
//
// Example usage:
// 		buffer := NewBufferMock()
// 		logger := NewZeroLog(buffer, 0)
// 		runner := &Runner{log: logger}
// 		go runner.run()
// 		buffer.WaitFor("finished running")
type BufferMock struct {
	bytes.Buffer
	t *testing.T

	notify  chan interface{}
	mNotify sync.Mutex
}

// NewBufferMock creates a new BufferMock
func NewBufferMock(t *testing.T) *BufferMock {
	return &BufferMock{
		t:      t,
		notify: nil,
	}
}

func (nb *BufferMock) Write(p []byte) (n int, err error) {
	n, err = nb.Buffer.Write(p)

	nb.mNotify.Lock()
	shouldNotify := nb.notify != nil
	nb.mNotify.Unlock()

	if shouldNotify {
		nb.notify <- struct{}{}
	}

	return n, err
}

// WaitFor waits for a target string to be found in the logs
// If log is not found, it timeouts after 5s without receiving logs
func (nb *BufferMock) WaitFor(target string) {
	nb.mNotify.Lock()
	nb.notify = make(chan interface{})
	nb.mNotify.Unlock()

	for {
		if strings.Contains(nb.String(), target) {
			break
		}
		select {
		case <-time.After(5 * time.Second):
			nb.t.Errorf("timeout after waiting for log: '%s'", target)
			nb.t.FailNow()
		case <-nb.notify:
		}
	}

	nb.mNotify.Lock()
	nb.notify = nil
	nb.mNotify.Unlock()
}
