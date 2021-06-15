package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	cb "github.com/minhdanh/ctrlcb/pkg/clipboard"
	"github.com/otiai10/copy"
)

func ProcessClipboardContent(clipboardContent, currentWorkingDirectory string, keepSourcePath, overwrite bool) int {
	processed := 0
	checkedFirstLine := false
	for _, line := range strings.Split(strings.TrimRight(clipboardContent, "\n"), "\n") {
		if !checkedFirstLine {
			if line != cb.ClipboardMarker {
				log.Printf("No file or directory paths found in clipboard. Doing nothing")
				return 0
			}
			checkedFirstLine = true
			continue
		}

		item := strings.Split(line, ":")
		sourceWorkingDirectory, err := cb.DecodeBase64(item[0])
		if err != nil {
			log.Printf("Cannot decode base64 string. Error: %v", err.Error())
			continue
		}
		sourcePath, err := cb.DecodeBase64(item[1])
		if err != nil {
			log.Printf("Cannot decode base64 string. Error: %v", err.Error())
			continue
		}

		err = CopyFileOrDir(sourceWorkingDirectory, sourcePath, currentWorkingDirectory, keepSourcePath, overwrite)
		if err != nil {
			log.Print(err.Error())
			continue
		}
		processed++
	}
	return processed
}

func CopyFileOrDir(sourceWorkingDirectory, sourcePath, currentWorkingDirectory string, keepSourcePath, overwrite bool) (err error) {
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
