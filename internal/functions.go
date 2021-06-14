package internal

import (
	"encoding/base64"
	"log"
	"strings"
)

func AddClipboardItem(clipboard_content, cwd, path string) (string, bool) {
	cwd_base64 := base64.StdEncoding.EncodeToString([]byte(cwd))
	path_base64 := base64.StdEncoding.EncodeToString([]byte(path))

	clipboard_item := cwd_base64 + ":" + path_base64
	if !strings.Contains(clipboard_content, clipboard_item) {
		clipboard_content = clipboard_content + "\n" + clipboard_item
		return clipboard_content, true
	}
	return clipboard_content, false
}
func DecodeBase64(str string) (string, error) {
	result_bytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Printf("Cannot decode base64 string. Error: %v", err.Error())
		return "", err
	}
	return string(result_bytes), nil
}
