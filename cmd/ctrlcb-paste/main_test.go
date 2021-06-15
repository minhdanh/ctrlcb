package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	functions "github.com/minhdanh/ctrlcb/internal"
)

// var keepSourcePath = flag.Bool("k", false, "keep the relative source path - create directories on destination directory if needed")

// // removeSource = flag.Bool("x", false, "remove the source file or directory")
// var overwrite = flag.Bool("f", false, "overwrite the destination file or directory if it already exists")

// func TestMain(m *testing.M) {
// 	flag.Parse()
// 	code := m.Run()
// 	// teardown(m)
// 	os.Exit(code)
// }

// func init() {
// 	// clipboard.Debug = true
// }

// func teardown(m *testing.M) {
// 	// os.RemoveAll("test/data/case03/case01")
// 	// os.RemoveAll("test/data.copy")
// 	// os.RemoveAll("test/data.copyTime")
// }

func TestInvalidClipboardContent(t *testing.T) {
	currentWorkingDirectory := t.TempDir()
	expected := 0
	clipboardContent := "invalid"

	actual := processItems(clipboardContent, currentWorkingDirectory, false, false)
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

	// should copy abs path without -k
	expected := 1
	clipboardContent := "# ctrlcb"
	clipboardContent, _ = functions.AddClipboardItem(clipboardContent, sourceWorkingDirectory, tmpFile.Name())

	actual := processItems(clipboardContent, targetWorkingDirectory, false, false)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should not copy abs path with -k
	expected = 0
	actual = processItems(clipboardContent, targetWorkingDirectory, true, false)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should not overwrite file
	expected = 0
	actual = processItems(clipboardContent, targetWorkingDirectory, false, false)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should overwrite file
	expected = 1
	actual = processItems(clipboardContent, targetWorkingDirectory, false, true)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}
}

func TestCopyFile_RelativePathTarget(t *testing.T) {
	sourceWorkingDirectory := t.TempDir()
	// sourcePath := filepath.Join(sourceWorkingDirectory, "a/b/c")
	// os.MkdirAll(sourcePath, os.ModePerm)
	tmpFileT, err := ioutil.TempFile(sourceWorkingDirectory, "test-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}

	targetWorkingDirectory := t.TempDir()
	os.Chdir(targetWorkingDirectory)
	// relativeFilePathT is relative from target directory
	relativeFilePathT, err := filepath.Rel(targetWorkingDirectory, tmpFileT.Name())
	if err != nil {
		log.Fatal("Cannot create relative path", err)
	}

	clipboardContent := "# ctrlcb"
	clipboardContent, _ = functions.AddClipboardItem(clipboardContent, sourceWorkingDirectory, relativeFilePathT)

	expected := 1
	actual := processItems(clipboardContent, targetWorkingDirectory, false, false)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should not overwrite file created in previous assertion
	expected = 0
	actual = processItems(clipboardContent, targetWorkingDirectory, false, false)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should overwrite file created in previous assertion
	expected = 1
	actual = processItems(clipboardContent, targetWorkingDirectory, false, true)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should copy with source path with -k
	expected = 1
	actual = processItems(clipboardContent, targetWorkingDirectory, true, false)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should not overwrite file created in previous assertion
	expected = 0
	actual = processItems(clipboardContent, targetWorkingDirectory, true, false)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should overwrite file created in previous assertion
	expected = 1
	actual = processItems(clipboardContent, targetWorkingDirectory, true, true)
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
	// sourcePath := filepath.Join(sourceWorkingDirectory, "a/b/c")
	// os.MkdirAll(sourcePath, os.ModePerm)
	tmpFileS, err := ioutil.TempFile(sourceWorkingDirectory, "test-")
	if err != nil {
		log.Fatal("Cannot create temporary file", err)
	}

	// relativeFilePathS is relative from source directory
	relativeFilePathS, err := filepath.Rel(sourceWorkingDirectory, tmpFileS.Name())
	if err != nil {
		log.Fatal("Cannot create relative path", err)
	}
	log.Print(relativeFilePathS)

	targetWorkingDirectory := t.TempDir()
	os.Chdir(targetWorkingDirectory)

	clipboardContent := "# ctrlcb"
	clipboardContent, _ = functions.AddClipboardItem(clipboardContent, sourceWorkingDirectory, relativeFilePathS)

	expected := 1
	actual := processItems(clipboardContent, targetWorkingDirectory, false, false)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should not overwrite file created in previous assertion
	expected = 0
	actual = processItems(clipboardContent, targetWorkingDirectory, false, false)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should overwrite file created in previous assertion
	expected = 1
	actual = processItems(clipboardContent, targetWorkingDirectory, false, true)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should not copy with source path with -k as there's no directory in path
	expected = 0
	actual = processItems(clipboardContent, targetWorkingDirectory, true, false)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// should overwrite file created in previous assertion
	expected = 1
	actual = processItems(clipboardContent, targetWorkingDirectory, true, true)
	if actual != expected {
		t.Fatalf("want %v, got %v", expected, actual)
	}

	// check if relative path exist in targetWorkingDirectory
	_, err = os.Stat(filepath.Join(targetWorkingDirectory, strings.TrimLeft(relativeFilePathS, "./")))
	if os.IsNotExist(err) {
		t.Fatalf("%v not copied to %v", relativeFilePathS, targetWorkingDirectory)
	}
}
