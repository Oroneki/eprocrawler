package main

import (
	"os"
	"testing"
)

func TestGetFolderInfo(t *testing.T) {
	setUpLoggers(os.Stderr, os.Stdout)
	folder := "C:\\Users\\00104509392\\Documents\\triagem\\"
	r, e := getFolderInfo(folder)
	if e != nil {
		t.Fatalf("e %v", e)
	}
	t.Logf("%#v %#v", r, len(r))
}

func TestAnotherGetFolderInfo(t *testing.T) {
	setUpLoggers(os.Stderr, os.Stdout)
	folder := "."
	r, e := getFolderInfo(folder)
	if e != nil {
		t.Fatalf("e %v", e)
	}
	t.Logf("%#v %#v", r, len(r))
}
