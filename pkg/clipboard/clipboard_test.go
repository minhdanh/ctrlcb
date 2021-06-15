package clipboard

import (
	"strings"
	"testing"
)

func TestDecodeBase64(t *testing.T) {
	base64Str := "MEVERjgzNTY3MUMwQzg2MjM4RUVFNEExRjREMTEzMzg="
	actual, err := DecodeBase64(base64Str)
	if err != nil {
		t.Fatal(err)
	}

	expected := "0EDF835671C0C86238EEE4A1F4D11338"
	if actual != expected {
		t.Fail()
	}

}

func TestAddClipboardItem(t *testing.T) {
	path := "a/b/c"
	cwd := "/tmp/test"

	_, added, _ := AddClipboardItem("", cwd, path)
	if added {
		t.Fatalf("should not add item without marker")
	}

	clipboardContent := ClipboardMarker
	clipboardContent, _, _ = AddClipboardItem(clipboardContent, cwd, path)

	expected := "L3RtcC90ZXN0:YS9iL2M="
	if !strings.Contains(clipboardContent, expected) {
		t.Fatalf("%v should contain %v", clipboardContent, expected)
	}

	_, added, _ = AddClipboardItem(clipboardContent, cwd, path)
	if added {
		t.Fatalf("should not add duplicated item")
	}
}
