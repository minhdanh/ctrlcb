package clipboard

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"
)

var ClipboardMarker = "# ctrlcb"

func AddClipboardItem(clipboardContent, cwd, path string) (string, bool, error) {
	if !strings.HasPrefix(clipboardContent, ClipboardMarker) {
		return clipboardContent, false, fmt.Errorf("cannot add item to clipboard as no marker detected")
	}

	cwdBase64 := base64.StdEncoding.EncodeToString([]byte(cwd))
	pathBase64 := base64.StdEncoding.EncodeToString([]byte(path))

	clipboardItem := cwdBase64 + ":" + pathBase64
	if !strings.Contains(clipboardContent, clipboardItem) {
		clipboardContent = clipboardContent + "\n" + clipboardItem
		return clipboardContent, true, nil
	}
	return clipboardContent, false, nil
}

func DecodeBase64(str string) (string, error) {
	result, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Printf("Cannot decode base64 string. Error: %v", err.Error())
		return "", err
	}
	return string(result), nil
}
