package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	functions "github.com/minhdanh/ctrlcb/internal"
)

func main() {
	args := os.Args[1:]
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current working directory: %s\n", cwd)

	clipboardContent := "# ctrlcb"

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
			log.Printf("File or directory %v does not exist.", path)
			continue
		}

		itemAdded := false
		clipboardContent, itemAdded = functions.AddClipboardItem(clipboardContent, cwd, args[i])
		if itemAdded {
			copiedPaths++
		}
	}

	if copiedPaths > 0 {
		fmt.Printf("Copied %v paths to clipboard\n", copiedPaths)
		clipboard.WriteAll(clipboardContent)
	}
}
