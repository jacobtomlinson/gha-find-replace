package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestHappyPath(t *testing.T) {
	_ = os.Mkdir("out", 0755)

	if err := ioutil.WriteFile("out/data.txt", []byte("Hello world\n"), 0644); err != nil {
		t.Error(err)
	}

	os.Setenv("INPUT_INCLUDE", "out/data.txt")
	os.Setenv("INPUT_EXCLUDE", "")
	os.Setenv("INPUT_FIND", "world")
	os.Setenv("INPUT_REPLACE", "there")
	os.Setenv("GITHUB_OUTPUT", "out/output.txt")

	main()

	data, err := ioutil.ReadFile("out/data.txt")
	if err != nil {
		t.Error(err)
	}

	want := "Hello there\n"
	if string(data) != want {
		t.Errorf("got %q, wanted %q", data, want)
	}
}
