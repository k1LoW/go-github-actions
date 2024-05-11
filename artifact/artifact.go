package artifact

import (
	"context"
	"io"

	"github.com/k1LoW/go-github-actions/artifact/legacy"
)

// Upload content as GitHub Actions artifact
func Upload(ctx context.Context, name, fp string, content io.Reader) error {
	return legacy.Upload(ctx, name, fp, content)
}

// UploadFiles as GitHub Actions artifact
func UploadFiles(ctx context.Context, name string, files []string) error {
	return legacy.UploadFiles(ctx, name, files)
}
