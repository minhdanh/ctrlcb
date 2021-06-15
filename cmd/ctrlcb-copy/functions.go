package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	cb "github.com/minhdanh/ctrlcb/pkg/clipboard"
)

func PrepareClipboardContent(cwd string, args []string) (string, int) {
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
	return clipboardContent, copiedPaths
}
