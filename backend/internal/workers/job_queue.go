package workers

import (
	"errors"
	"sync"
)

var ErrQueueFull = errors.New("submission queue is full")

type JobQueue interface {
	Enqueue(submissionID string) error
	Channel() <-chan string
	Close()
}

type ChannelJobQueue struct {
	mu     sync.RWMutex
	ch     chan string
	closed bool
}

func NewChannelJobQueue(bufferSize int) *ChannelJobQueue {
	if bufferSize <= 0 {
		bufferSize = 1000
	}
	return &ChannelJobQueue{
		ch: make(chan string, bufferSize),
	}
}

func (q *ChannelJobQueue) Enqueue(submissionID string) error {
	q.mu.RLock()
	defer q.mu.RUnlock()
	if q.closed {
		return errors.New("queue is closed")
	}
	select {
	case q.ch <- submissionID:
		return nil
	default:
		return ErrQueueFull
	}
}

func (q *ChannelJobQueue) Channel() <-chan string {
	return q.ch
}

func (q *ChannelJobQueue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	if !q.closed {
		q.closed = true
		close(q.ch)
	}
}
