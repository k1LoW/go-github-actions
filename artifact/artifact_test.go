package artifact

import (
	"os"
	"strings"
	"testing"
)

func TestUpload(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "" {
		t.Skip("Not running on GitHub Actions")
	}
	if err := Upload(t.Context(), "TestUpload", "artifact/testdata/test.txt", strings.NewReader("hello artifact 3\n")); err != nil {
		t.Error(err)
	}
}

func TestUploadLargeContent(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "" {
		t.Skip("Not running on GitHub Actions")
	}
	s := strings.Repeat("0123456789\n", 1024*1024*10)
	name := "TestUploadLargeContent"
	if err := Upload(t.Context(), name, "artifact/testdata/large.txt", strings.NewReader(s)); err != nil {
		t.Error(err)
	}
}

func TestUploadFiles(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "" {
		t.Skip("Not running on GitHub Actions")
	}
	files := []string{
		"testdata/test2.txt",
		"testdata/test3.txt",
	}
	if err := UploadFiles(t.Context(), "TestUploadFiles", files); err != nil {
		t.Error(err)
	}
}
