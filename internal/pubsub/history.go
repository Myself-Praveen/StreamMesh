package pubsub

import (
	"sync"
)

// HistoryBuffer is a thread-safe ring buffer for storing recent messages
type HistoryBuffer struct {
	messages [][]byte
	head     int
	tail     int
	count    int
	capacity int
	mu       sync.RWMutex
}

// NewHistoryBuffer creates a new HistoryBuffer with the given capacity
func NewHistoryBuffer(capacity int) *HistoryBuffer {
	return &HistoryBuffer{
		messages: make([][]byte, capacity),
		capacity: capacity,
	}
}

// Add appends a new message to the buffer, overwriting the oldest if full
func (b *HistoryBuffer) Add(msg []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.capacity == 0 {
		return
	}

	b.messages[b.tail] = msg
	b.tail = (b.tail + 1) % b.capacity

	if b.count < b.capacity {
		b.count++
	} else {
		b.head = (b.head + 1) % b.capacity
	}
}

// GetRecent returns up to n most recent messages from the buffer
func (b *HistoryBuffer) GetRecent(n int) [][]byte {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.count == 0 || n <= 0 {
		return nil
	}

	if n > b.count {
		n = b.count
	}

	result := make([][]byte, n)
	
	// Start reading from the most recent (tail - 1) backwards
	// To make the output chronological, we read forwards from (tail - n)
	
	startIndex := (b.tail - n)
	if startIndex < 0 {
		startIndex += b.capacity
	}

	for i := 0; i < n; i++ {
		idx := (startIndex + i) % b.capacity
		result[i] = b.messages[idx]
	}

	return result
}
