package ws

import (
	"sync"
	"testing"
)

func TestManager_AddRemoveGet(t *testing.T) {
	manager := NewManager()

	conn1 := NewConnection("c1", nil)
	conn2 := NewConnection("c2", nil)

	manager.Add(conn1)
	manager.Add(conn2)

	if manager.Count() != 2 {
		t.Errorf("Expected 2 connections, got %d", manager.Count())
	}

	c, ok := manager.Get("c1")
	if !ok || c.ID != "c1" {
		t.Errorf("Expected to get connection c1")
	}

	manager.Remove("c1")
	if manager.Count() != 1 {
		t.Errorf("Expected 1 connection after removal, got %d", manager.Count())
	}

	_, ok = manager.Get("c1")
	if ok {
		t.Errorf("Expected connection c1 to be removed")
	}
}

func TestManager_Concurrency(t *testing.T) {
	manager := NewManager()
	var wg sync.WaitGroup

	// Concurrently add 1000 connections
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			conn := NewConnection(id, nil)
			manager.Add(conn)
		}(string(rune(i)))
	}

	wg.Wait()

	if manager.Count() != 1000 {
		t.Errorf("Expected 1000 connections, got %d", manager.Count())
	}

	// Concurrently remove them
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			manager.Remove(id)
		}(string(rune(i)))
	}

	wg.Wait()

	if manager.Count() != 0 {
		t.Errorf("Expected 0 connections, got %d", manager.Count())
	}
}
