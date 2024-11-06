package main

import (
	"os"
	"testing"
)

func TestHappyPath(t *testing.T) {
	_ = os.Mkdir("out", 0o755)

	if err := os.WriteFile("out/data.txt", []byte("Hello world\n"), 0o644); err != nil {
		t.Error(err)
	}

	t.Setenv("INPUT_INCLUDE", "out/data.txt")
	t.Setenv("INPUT_EXCLUDE", "")
	t.Setenv("INPUT_FIND", "world")
	t.Setenv("INPUT_REPLACE", "there")
	t.Setenv("GITHUB_OUTPUT", "out/output.txt")

	main()

	data, err := os.ReadFile("out/data.txt")
	if err != nil {
		t.Error(err)
	}

	want := "Hello there\n"
	if string(data) != want {
		t.Errorf("got %q, wanted %q", data, want)
	}
}
