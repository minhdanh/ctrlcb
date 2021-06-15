package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/atotto/clipboard"
)

func main() {
	keepSourcePath := flag.Bool("k", false, "keep the relative source path - create directories on destination directory if needed")
	overwrite := flag.Bool("f", false, "overwrite the destination file or directory if it already exists")
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

	processedItems := ProcessClipboardContent(clipboardContent, currentWorkingDirectory, *keepSourcePath, *overwrite)
	log.Printf("%v items were copied", processedItems)
}
