package artifact

import (
	"context"
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
