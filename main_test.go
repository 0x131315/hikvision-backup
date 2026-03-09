package main

import (
	"reflect"
	"testing"
)

func TestNormalizeLegacyArgs(t *testing.T) {
	in := []string{"-vv", "-vvv", "-v", "--other"}
	got := normalizeLegacyArgs(in)
	want := []string{"--verbose", "--verbose-http", "-v", "--other"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected normalized args: got=%v want=%v", got, want)
	}
}
