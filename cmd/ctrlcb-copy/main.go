package main

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
)

func main() {
	args := os.Args[1:]
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	fmt.Printf("Current working directory: %s\n", cwd)
	clipboardContent, copiedPaths := PrepareClipboardContent(cwd, args)

	if copiedPaths > 0 {
		fmt.Printf("Copied %v paths to clipboard\n", copiedPaths)
		clipboard.WriteAll(clipboardContent)
	}
}
