package main_test

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/minhdanh/ctrlcb/cmd/ctrlcb-paste"
	cb "github.com/minhdanh/ctrlcb/pkg/clipboard"
)

func TestInvalidClipboardContent(t *testing.T) {
	currentWorkingDirectory := t.TempDir()
	expected := 0
	clipboardContent := "invalid"

	actual := ProcessClipboardContent(clipboardContent, currentWorkingDirectory, false, false)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}
}

func TestCopyFile_AbsolutePath(t *testing.T) {
	sourceWorkingDirectory := t.TempDir()
	tmpFile, err := ioutil.TempFile(sourceWorkingDirectory, "test-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}

	targetWorkingDirectory := t.TempDir()
	os.Chdir(targetWorkingDirectory)

	// should copy absolute path without -k flag
	expected := 1
	clipboardContent := cb.ClipboardMarker
	clipboardContent, _, _ = cb.AddClipboardItem(clipboardContent, sourceWorkingDirectory, tmpFile.Name())

	actual := ProcessClipboardContent(clipboardContent, targetWorkingDirectory, false, false)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should not copy absolute path with -k flag
	expected = 0
	actual = ProcessClipboardContent(clipboardContent, targetWorkingDirectory, true, false)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should not overwrite file without -f flag
	expected = 0
	actual = ProcessClipboardContent(clipboardContent, targetWorkingDirectory, false, false)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should overwrite file with -f flag
	expected = 1
	actual = ProcessClipboardContent(clipboardContent, targetWorkingDirectory, false, true)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}
}

func TestCopyFile_RelativePathTarget(t *testing.T) {
	sourceWorkingDirectory := t.TempDir()
	tmpFileT, err := ioutil.TempFile(sourceWorkingDirectory, "test-")
	if err != nil {
		t.Fatal("Cannot create temporary file", err)
	}

	targetWorkingDirectory := t.TempDir()
	os.Chdir(targetWorkingDirectory)

	// relativeFilePathT is relative from target directory
	relativeFilePathT, err := filepath.Rel(targetWorkingDirectory, tmpFileT.Name())
	if err != nil {
		t.Fatal("Cannot create relative path", err)
	}

	clipboardContent := cb.ClipboardMarker
	clipboardContent, _, _ = cb.AddClipboardItem(clipboardContent, sourceWorkingDirectory, relativeFilePathT)

	expected := 1
	actual := ProcessClipboardContent(clipboardContent, targetWorkingDirectory, false, false)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should not overwrite file created in previous assertion
	expected = 0
	actual = ProcessClipboardContent(clipboardContent, targetWorkingDirectory, false, false)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should overwrite file created in previous assertion
	expected = 1
	actual = ProcessClipboardContent(clipboardContent, targetWorkingDirectory, false, true)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should copy with source path with -k flag
	expected = 1
	actual = ProcessClipboardContent(clipboardContent, targetWorkingDirectory, true, false)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should not overwrite file created in previous assertion
	expected = 0
	actual = ProcessClipboardContent(clipboardContent, targetWorkingDirectory, true, false)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should overwrite file created in previous assertion
	expected = 1
	actual = ProcessClipboardContent(clipboardContent, targetWorkingDirectory, true, true)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// check if relative path exist in targetWorkingDirectory
	_, err = os.Stat(filepath.Join(targetWorkingDirectory, strings.TrimLeft(relativeFilePathT, "./")))
	if os.IsNotExist(err) {
		t.Fatalf("%v not copied to %v", relativeFilePathT, targetWorkingDirectory)
	}
}

func TestCopyFile_RelativePathSource(t *testing.T) {
	sourceWorkingDirectory := t.TempDir()
	tmpFileS, err := ioutil.TempFile(sourceWorkingDirectory, "test-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}

	// relativeFilePathS is relative from source directory
	relativeFilePathS, err := filepath.Rel(sourceWorkingDirectory, tmpFileS.Name())
	if err != nil {
		t.Fatal("Cannot create relative path", err)
	}

	targetWorkingDirectory := t.TempDir()
	os.Chdir(targetWorkingDirectory)

	clipboardContent := cb.ClipboardMarker
	clipboardContent, _, _ = cb.AddClipboardItem(clipboardContent, sourceWorkingDirectory, relativeFilePathS)

	expected := 1
	actual := ProcessClipboardContent(clipboardContent, targetWorkingDirectory, false, false)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should not overwrite file created in previous assertion
	expected = 0
	actual = ProcessClipboardContent(clipboardContent, targetWorkingDirectory, false, false)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should overwrite file created in previous assertion
	expected = 1
	actual = ProcessClipboardContent(clipboardContent, targetWorkingDirectory, false, true)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should not copy with source path with -k flag as path is just a filename
	expected = 0
	actual = ProcessClipboardContent(clipboardContent, targetWorkingDirectory, true, false)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should overwrite file created in previous assertion
	expected = 1
	actual = ProcessClipboardContent(clipboardContent, targetWorkingDirectory, true, true)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// check if relative path exist in targetWorkingDirectory
	_, err = os.Stat(filepath.Join(targetWorkingDirectory, strings.TrimLeft(relativeFilePathS, "./")))
	if os.IsNotExist(err) {
		t.Fatalf("%v not copied to %v", relativeFilePathS, targetWorkingDirectory)
	}
}
