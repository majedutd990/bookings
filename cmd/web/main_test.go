package main

import "testing"

//TestRun tests run function in main
func TestRun(t *testing.T) {
	err := run()
	if err != nil {
		t.Error("failed run()")
	}
}
