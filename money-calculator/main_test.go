package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func Test_main(t *testing.T) {

	stdOut := os.Stdout

	r, w, _ := os.Pipe()

	os.Stdout = w

	main()

	_ = w.Close()

	result, _ := ioutil.ReadAll(r)

	output := string(result)

	os.Stdout = stdOut

	expected := "Final bank balance"

	if !strings.Contains(output, expected) {

		t.Errorf("Expected %s, got %s", expected, output)

	}

}
