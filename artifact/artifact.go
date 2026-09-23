package artifact

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"connectrpc.com/connect"
	"github.com/k1LoW/go-github-actions/artifact/legacy"
	apiv1 "github.com/k1LoW/go-github-actions/artifact/proto/gen/go/results/api/v1"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Upload content as GitHub Actions artifact.
func Upload(ctx context.Context, name, fp string, content io.Reader, opts ...Option) error {
	c := newConfig(opts)
	if useLegacy() {
		if err := c.handleAttestError(errors.New("attestation is not supported with legacy artifact upload")); err != nil {
			return err
		}
		return legacy.Upload(ctx, name, fp, content)
	}

	return uploadZip(ctx, c, name, func(zw *zip.Writer) error {
		h := &zip.FileHeader{
			Name:   fp,
			Method: zip.Deflate,
		}
		w, err := zw.CreateHeader(h)
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, content); err != nil {
			return err
		}
		return nil
	})
}

// UploadFiles as GitHub Actions artifact.
func UploadFiles(ctx context.Context, name string, files []string, opts ...Option) error {
	c := newConfig(opts)
	if useLegacy() {
		if err := c.handleAttestError(errors.New("attestation is not supported with legacy artifact upload")); err != nil {
			return err
		}
		return legacy.UploadFiles(ctx, name, files)
	}

	return uploadZip(ctx, c, name, func(zw *zip.Writer) error {
		for _, fp := range files {
			if err := func() error {
				a, err := filepath.Abs(fp)
				if err != nil {
					return err
				}
				rel, err := filepath.Rel(os.Getenv("GITHUB_WORKSPACE"), a)
				if err != nil {
					return err
				}
				f, err := os.Open(fp)
				if err != nil {
					return err
				}
				defer f.Close()
				fi, err := f.Stat()
				if err != nil {
					return err
				}
				h, err := zip.FileInfoHeader(fi)
				if err != nil {
					return err
				}
				h.Name = rel
				h.Method = zip.Deflate
				w, err := zw.CreateHeader(h)
				if err != nil {
					return err
				}
				if _, err := io.Copy(w, f); err != nil {
					return err
				}
				return nil
			}(); err != nil {
				return err
			}
		}
		return nil
	})
}

// ErrUnarchivedUploadNotSupported is returned by UploadUnarchived where artifacts are
// uploaded through the legacy API, which archives every artifact it serves.
var ErrUnarchivedUploadNotSupported = errors.New("uploading an artifact without archiving it is not supported with legacy artifact upload")

// UploadUnarchived uploads content as a GitHub Actions artifact of a single file, without
// archiving it, and returns the ID of the artifact. The artifact is named after the file, as
// actions/upload-artifact names it with `archive: false`, and is served as the file itself,
// so an HTML file opens in the browser rather than being downloaded as a zip.
func UploadUnarchived(ctx context.Context, name string, content io.Reader, opts ...Option) (int64, error) {
	c := newConfig(opts)
	if useLegacy() {
		return 0, ErrUnarchivedUploadNotSupported
	}
	b, err := io.ReadAll(content)
	if err != nil {
		return 0, err
	}
	return uploadBlob(ctx, c, name, mimeType(name), b)
}

func uploadZip(ctx context.Context, c *config, name string, write func(zw *zip.Writer) error) error {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	if err := write(zw); err != nil {
		return err
	}
	if err := zw.Close(); err != nil {
		return err
	}
	_, err := uploadBlob(ctx, c, name, "", buf.Bytes())
	return err
}

// uploadBlob creates the artifact, uploads b as its content and finalizes it. An empty
// contentType uploads b as the zip archive the artifact is served as.
func uploadBlob(ctx context.Context, c *config, name, contentType string, b []byte) (int64, error) {
	ids, err := getBackendIdsFromToken()
	if err != nil {
		return 0, err
	}
	apic, err := newAPIClient()
	if err != nil {
		return 0, err
	}
	creq := &apiv1.CreateArtifactRequest{
		WorkflowRunBackendId:    ids.workflowRunBackendId,
		WorkflowJobRunBackendId: ids.workflowJobRunBackendId,
		Name:                    name,
		Version:                 4,
	}
	if contentType != "" {
		// Version 7 is the one that reads the MIME type, which is what tells the backend the
		// content is the file itself rather than a zip archive of it.
		creq.Version = 7
		creq.MimeType = wrapperspb.String(contentType)
	}

	res, err := apic.CreateArtifact(ctx, connect.NewRequest(creq))
	if err != nil {
		return 0, err
	}
	if !res.Msg.GetOk() {
		return 0, errors.New("response is not ok")
	}

	size := int64(len(b))
	sum := sha256.Sum256(b)
	digest := hex.EncodeToString(sum[:])
	if err := upload(ctx, res.Msg.GetSignedUploadUrl(), bytes.NewReader(b), contentType); err != nil {
		return 0, err
	}

	// Attest before finalizing so that a failed required attestation does not leave a visible, unattested artifact.
	if err := c.attest(ctx, name, map[string]string{"sha256": digest}); err != nil {
		return 0, err
	}

	{
		req := connect.NewRequest(&apiv1.FinalizeArtifactRequest{
			WorkflowRunBackendId:    ids.workflowRunBackendId,
			WorkflowJobRunBackendId: ids.workflowJobRunBackendId,
			Name:                    name,
			Size:                    size,
			Hash:                    wrapperspb.String("sha256:" + digest),
		})

		res, err := apic.FinalizeArtifact(ctx, req)
		if err != nil {
			return 0, err
		}
		if !res.Msg.GetOk() {
			return 0, errors.New("response is not ok")
		}
		return res.Msg.GetArtifactId(), nil
	}
}

// mimeType returns the MIME type an artifact named after a file is served with.
func mimeType(name string) string {
	if t := mime.TypeByExtension(filepath.Ext(name)); t != "" {
		return t
	}
	return "application/octet-stream"
}

func useLegacy() bool {
	if isGHES() {
		return true
	}
	if os.Getenv("ACTIONS_USE_LEGACY_ARTIFACT_UPLOAD") != "" {
		return true
	}
	return false
}

func isGHES() bool {
	if os.Getenv("GITHUB_SERVER_URL") == "" {
		return false
	}
	if strings.HasSuffix(os.Getenv("GITHUB_SERVER_URL"), "github.com") {
		return false
	}
	if strings.HasSuffix(os.Getenv("GITHUB_SERVER_URL"), ".ghe.com") {
		return false
	}
	if strings.HasSuffix(os.Getenv("GITHUB_SERVER_URL"), ".ghe.localhost") {
		return false
	}
	return true
}

type backendIds struct {
	workflowRunBackendId    string
	workflowJobRunBackendId string
}

func getBackendIdsFromToken() (*backendIds, error) {
	rt := os.Getenv("ACTIONS_RUNTIME_TOKEN")
	if rt == "" {
		return nil, errors.New("env ACTIONS_RUNTIME_TOKEN is only available from the context of an action")
	}
	jt, err := jwt.ParseString(rt, jwt.WithVerify(false), jwt.WithValidate(false))
	if err != nil {
		return nil, err
	}
	scp, ok := jt.Get("scp")
	if !ok {
		return nil, errors.New("no scp in ACTIONS_RUNTIME_TOKEN")
	}
	scpParts, ok := scp.(string)
	if !ok {
		return nil, errors.New("invalid scp in ACTIONS_RUNTIME_TOKEN")
	}
	for _, scopes := range strings.Split(scpParts, " ") {
		scopeParts := strings.Split(scopes, ":")
		if scopeParts[0] != "Actions.Results" {
			continue
		}
		if len(scopeParts) != 3 {
			return nil, errors.New("invalid scp in ACTIONS_RUNTIME_TOKEN")
		}
		return &backendIds{
			workflowRunBackendId:    scopeParts[1],
			workflowJobRunBackendId: scopeParts[2],
		}, nil
	}

	return nil, errors.New("no backend ids found")
}
