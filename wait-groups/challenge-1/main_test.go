package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func Test_updateMessage(t *testing.T) {

	wg.Add(1)

	go updateMessage("This is a string!", &wg)

	wg.Wait()

	if msg != "This is a string!" {
		t.Errorf("updateMessage() = %v, want %v", msg, "This is a string!")
	}
}

func Test_printMessage(t *testing.T) {
	stdOut := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	msg = "This is a string!"
	printMessage()

	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	if !strings.Contains(output, "This is a string!") {
		t.Errorf("printMessage() = %v, want %v", output, "This is a string!")
	}
}

func Test_main(t *testing.T) {
	stdOut := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	if !strings.Contains(output, "Hello, universe!") {
		t.Errorf("main() = %v, want %v", output, "Hello, universe!")
	}

	if !strings.Contains(output, "Hello, cosmos!") {
		t.Errorf("main() = %v, want %v", output, "Hello, cosmos!")
	}

	if !strings.Contains(output, "Hello, world!") {
		t.Errorf("main() = %v, want %v", output, "Hello, world!")
	}
}
