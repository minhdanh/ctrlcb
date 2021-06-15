package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	cb "github.com/minhdanh/ctrlcb/pkg/clipboard"
)

func main() {
	args := os.Args[1:]
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	fmt.Printf("Current working directory: %s\n", cwd)

	clipboardContent := cb.ClipboardMarker

	copiedPaths := 0
	for i := 0; i < len(args); i++ {
		usr, _ := user.Current()
		home_dir := usr.HomeDir

		path := args[i]
		if path == "~" {
			path = home_dir
		} else if strings.HasPrefix(path, "~/") {
			path = filepath.Join(home_dir, path[2:])
		}

		if !filepath.IsAbs(path) {
			path = filepath.Join(cwd, path)
		}
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			fmt.Printf("File or directory %v does not exist\n", path)
			continue
		}

		itemAdded := false
		clipboardContent, itemAdded, err = cb.AddClipboardItem(clipboardContent, cwd, args[i])
		if err != nil {
			panic(err)
		}
		if itemAdded {
			copiedPaths++
		}
	}

	if copiedPaths > 0 {
		fmt.Printf("Copied %v paths to clipboard\n", copiedPaths)
		clipboard.WriteAll(clipboardContent)
	}
}
