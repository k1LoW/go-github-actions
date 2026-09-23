package artifact

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/k1LoW/go-github-actions/attest"
)

func TestUpload(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "" {
		t.Skip("Not running on GitHub Actions")
	}
	if err := Upload(context.TODO(), "TestUpload", "artifact/testdata/test.txt", strings.NewReader("hello artifact 3\n")); err != nil {
		t.Error(err)
	}
}

func TestUploadLargeContent(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "" {
		t.Skip("Not running on GitHub Actions")
	}
	ctx := context.TODO()
	s := strings.Repeat("0123456789\n", 1024*1024*10)
	name := "TestUploadLargeContent"
	if err := Upload(ctx, name, "artifact/testdata/large.txt", strings.NewReader(s)); err != nil {
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
	if err := UploadFiles(context.TODO(), "TestUploadFiles", files); err != nil {
		t.Error(err)
	}
}

func TestUploadWithAttestation(t *testing.T) {
	if os.Getenv("ATTEST_E2E") == "" {
		t.Skip("ATTEST_E2E is not set")
	}
	// The CI downloads this artifact as a zip and verifies it with `gh attestation verify`.
	opt := WithAttestation(AttestRequired, attest.New())
	if err := Upload(context.TODO(), "TestUploadWithAttestation", "artifact/testdata/attested.txt", strings.NewReader("hello attestation\n"), opt); err != nil {
		t.Error(err)
	}
}

func TestUploadUnarchived(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "" {
		t.Skip("Not running on GitHub Actions")
	}
	id, err := UploadUnarchived(context.TODO(), "TestUploadUnarchived.html", strings.NewReader("<!doctype html><title>hello artifact</title>\n"))
	if err != nil {
		t.Fatal(err)
	}
	if id == 0 {
		t.Error("the ID of the uploaded artifact is not returned")
	}
}

func TestUploadUnarchivedOnLegacy(t *testing.T) {
	t.Setenv("ACTIONS_USE_LEGACY_ARTIFACT_UPLOAD", "true")
	if _, err := UploadUnarchived(context.TODO(), "test.html", strings.NewReader("")); !errors.Is(err, ErrUnarchivedUploadNotSupported) {
		t.Errorf("got %v, want %v", err, ErrUnarchivedUploadNotSupported)
	}
}

func TestMimeType(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"report.html", "text/html"},
		{"report.json", "application/json"},
		{"report", "application/octet-stream"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// A prefix, since the system MIME tables can add parameters such as a charset.
			if got := mimeType(tt.name); !strings.HasPrefix(got, tt.want) {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
