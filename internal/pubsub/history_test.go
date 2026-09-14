package pubsub

import (
	"reflect"
	"testing"
)

func TestHistoryBuffer(t *testing.T) {
	buffer := NewHistoryBuffer(3)

	buffer.Add([]byte("msg1"))
	buffer.Add([]byte("msg2"))

	recent := buffer.GetRecent(2)
	if len(recent) != 2 || string(recent[0]) != "msg1" || string(recent[1]) != "msg2" {
		t.Errorf("Expected msg1, msg2, got %v", recent)
	}

	buffer.Add([]byte("msg3"))
	buffer.Add([]byte("msg4")) // should overwrite msg1

	recent = buffer.GetRecent(3)
	if len(recent) != 3 || string(recent[0]) != "msg2" || string(recent[1]) != "msg3" || string(recent[2]) != "msg4" {
		t.Errorf("Expected msg2, msg3, msg4, got %v", recent)
	}
	
	// Test requesting more than capacity
	recent = buffer.GetRecent(5)
	if len(recent) != 3 {
		t.Errorf("Expected 3 items when requesting more than buffer size")
	}

	// Test string equality with reflection for easier debugging if it fails
	expected := [][]byte{[]byte("msg2"), []byte("msg3"), []byte("msg4")}
	if !reflect.DeepEqual(recent, expected) {
		t.Errorf("Expected %v, got %v", expected, recent)
	}
}
