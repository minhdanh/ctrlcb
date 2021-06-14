package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	functions "github.com/minhdanh/ctrlcb/internal"
	"github.com/otiai10/copy"
)

func main() {
	keepSourcePath := flag.Bool("k", false, "keep the relative source path - create directories on destination directory if needed")
	// removeSource = flag.Bool("x", false, "remove the source file or directory")
	overwrite := flag.Bool("f", false, "overwrite the destination file or directory if it already exists")
	// symlink
	// mode
	// owner
	flag.Parse()

	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current working directory: %s\n", currentWorkingDirectory)

	clipboardContent, err := clipboard.ReadAll()
	if err != nil {
		panic(err)
	}

	processedItems := processItems(clipboardContent, currentWorkingDirectory, *keepSourcePath, *overwrite)
	log.Printf("%v items were copied", processedItems)
}

func processItems(clipboardContent, currentWorkingDirectory string, keepSourcePath, overwrite bool) int {
	processed := 0
	checkedFirstLine := false
	for _, line := range strings.Split(strings.TrimRight(clipboardContent, "\n"), "\n") {
		if !checkedFirstLine {
			if line != "# ctrlcb" {
				log.Printf("No file or directory paths found in clipboard. Doing nothing")
				return 0
			}
			checkedFirstLine = true
			continue
		}

		item := strings.Split(line, ":")
		sourceWorkingDirectory, err := functions.DecodeBase64(item[0])
		if err != nil {
			log.Printf("Cannot decode base64 string. Error: %v", err.Error())
			continue
		}
		sourcePath, err := functions.DecodeBase64(item[1])
		if err != nil {
			log.Printf("Cannot decode base64 string. Error: %v", err.Error())
			continue
		}

		err = copyFileOrDir(sourceWorkingDirectory, sourcePath, currentWorkingDirectory, keepSourcePath, overwrite)
		if err != nil {
			log.Print(err.Error())
			continue
		}
		processed++
	}
	return processed
}

func copyFileOrDir(sourceWorkingDirectory, sourcePath, currentWorkingDirectory string, keepSourcePath, overwrite bool) (err error) {
	absSourcePath := sourcePath
	log.Printf("Path: %v", sourcePath)
	if filepath.IsAbs(sourcePath) {
		log.Printf("%v is an absolute path", sourcePath)
		if keepSourcePath {
			return fmt.Errorf("cannot use -k with absolute path %v", sourcePath)
		}
	} else {
		absSourcePath = filepath.Join(sourceWorkingDirectory, absSourcePath)
	}

	log.Printf("Path: %v", absSourcePath)
	_, err = os.Stat(absSourcePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("source file or directory %v does not exist", absSourcePath)
	}

	name := filepath.Base(sourcePath)
	target := filepath.Join(currentWorkingDirectory, name)

	if keepSourcePath {
		relativeTargetDir := filepath.Dir(strings.TrimLeft(sourcePath, "./"))
		os.MkdirAll(filepath.Join(currentWorkingDirectory, relativeTargetDir), os.ModePerm)
		target = filepath.Join(currentWorkingDirectory, relativeTargetDir, name)
	}

	_, err = os.Stat(target)
	if !os.IsNotExist(err) {
		if overwrite {
			log.Printf("Target already exists. Overwriting")
		} else {
			return errors.New("target already exists. To overwrite please use -f flag")
		}
	}

	opt := copy.Options{
		Skip: func(src string) (bool, error) {
			return strings.HasSuffix(src, ".git"), nil
		},
	}
	err = copy.Copy(absSourcePath, target, opt)
	if err != nil {
		return err
	}
	log.Printf("Copied %v to %v", absSourcePath, target)
	return nil
}
