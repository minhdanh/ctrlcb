// TODO:
// - test if function return readable clipboard content format?
// - test if function work with ~
package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
)

func TestCopyNonExistFiles(t *testing.T) {
	cwd := t.TempDir()
	args := []string{"file1", "file2", "file3"}

	expected := 0
	_, copiedPaths := PrepareClipboardContent(cwd, args)
	if copiedPaths != expected {
		t.Fatalf("want %v, got %v", expected, copiedPaths)
	}
}
func TestCopyDuplicatedPaths(t *testing.T) {
	cwd := t.TempDir()
	tmpFile, err := ioutil.TempFile(cwd, "test-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}

	relFileName, _ := filepath.Rel(cwd, tmpFile.Name())
	args := []string{relFileName, relFileName}

	expected := 1
	_, copiedPaths := PrepareClipboardContent(cwd, args)
	if copiedPaths != expected {
		t.Fatalf("want %v, got %v", expected, copiedPaths)
	}
}

func TestCopyAbsolutePaths(t *testing.T) {
	cwd := t.TempDir()
	tmpFile, err := ioutil.TempFile(cwd, "test-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}

	args := []string{tmpFile.Name()}

	expected := 1
	_, copiedPaths := PrepareClipboardContent(cwd, args)
	if copiedPaths != expected {
		t.Fatalf("want %v, got %v", expected, copiedPaths)
	}
}
