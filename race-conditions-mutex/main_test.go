package main

import (
	"sync"
	"testing"
)

func Test_updateMessage(t *testing.T) {

	msg = "Hello, World!"

	var mutex sync.Mutex

	wg.Add(1)
	go updateMessage("Hello, Go!", &mutex)
	wg.Wait()

	if msg != "Hello, Go!" {
		t.Errorf("msg = %s; want Hello, Go!", msg)
	}
}
