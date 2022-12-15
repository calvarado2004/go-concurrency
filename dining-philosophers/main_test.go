package main

import (
	"testing"
	"time"
)

func Test_dine(t *testing.T) {

	eatTime = 0 * time.Second
	thinkTime = 0 * time.Second
	sleepTime = 0 * time.Second

	for i := 0; i < 10; i++ {
		orderFinished = []string{}
		dine()
		if len(orderFinished) != len(philosophers) {
			t.Errorf("Expected %d philosophers to finish, got %d", len(philosophers), len(orderFinished))
		}
	}

}

func Test_dineWithVaryingDelays(t *testing.T) {

	var theTests = []struct {
		name  string
		delay time.Duration
	}{
		{"zero delay", 0 * time.Second},
		{"quarter second delay", 250 * time.Millisecond},
		{"half second delay", 500 * time.Millisecond},
	}

	for _, tt := range theTests {
		orderFinished = []string{}
		eatTime = tt.delay
		thinkTime = tt.delay
		sleepTime = tt.delay

		dine()
		if len(orderFinished) != len(philosophers) {
			t.Errorf("Expected %d philosophers to finish, got %d", len(philosophers), len(orderFinished))
		}
	}

}
