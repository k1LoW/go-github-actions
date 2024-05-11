package artifact

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"connectrpc.com/connect"
	"github.com/k1LoW/go-github-actions/artifact/legacy"
	apiv1 "github.com/k1LoW/go-github-actions/artifact/proto/gen/go/results/api/v1"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// Upload content as GitHub Actions artifact
func Upload(ctx context.Context, name, fp string, content io.Reader) error {
	if useLegacy() {
		return legacy.Upload(ctx, name, fp, content)
	}

	ids, err := getBackendIdsFromToken()
	if err != nil {
		return err
	}
	apic, err := newAPIClient()
	if err != nil {
		return err
	}
	req := connect.NewRequest(&apiv1.CreateArtifactRequest{
		WorkflowRunBackendId:    ids.workflowRunBackendId,
		WorkflowJobRunBackendId: ids.workflowJobRunBackendId,
		Name:                    name,
		Version:                 4,
	})

	res, err := apic.CreateArtifact(ctx, req)
	if err != nil {
		return err
	}
	if !res.Msg.GetOk() {
		return errors.New("response is not ok")
	}

	var size int64
	{
		buf := new(bytes.Buffer)
		zw := zip.NewWriter(buf)
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
		if err := zw.Close(); err != nil {
			return err
		}
		if err := upload(ctx, res.Msg.GetSignedUploadUrl(), buf); err != nil {
			return err
		}
		size = int64(buf.Len())
	}

	{
		req := connect.NewRequest(&apiv1.FinalizeArtifactRequest{
			WorkflowRunBackendId:    ids.workflowRunBackendId,
			WorkflowJobRunBackendId: ids.workflowJobRunBackendId,
			Name:                    name,
			Size:                    size,
		})

		res, err := apic.FinalizeArtifact(ctx, req)
		if err != nil {
			return err
		}
		if !res.Msg.GetOk() {
			return errors.New("response is not ok")
		}
	}

	return nil
}

// UploadFiles as GitHub Actions artifact
func UploadFiles(ctx context.Context, name string, files []string) error {
	if useLegacy() {
		return legacy.UploadFiles(ctx, name, files)
	}

	ids, err := getBackendIdsFromToken()
	if err != nil {
		return err
	}
	apic, err := newAPIClient()
	if err != nil {
		return err
	}
	req := connect.NewRequest(&apiv1.CreateArtifactRequest{
		WorkflowRunBackendId:    ids.workflowRunBackendId,
		WorkflowJobRunBackendId: ids.workflowJobRunBackendId,
		Name:                    name,
		Version:                 4,
	})

	res, err := apic.CreateArtifact(ctx, req)
	if err != nil {
		return err
	}
	if !res.Msg.GetOk() {
		return errors.New("response is not ok")
	}

	var size int64
	{
		buf := new(bytes.Buffer)
		zw := zip.NewWriter(buf)
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
		if err := zw.Close(); err != nil {
			return err
		}
		if err := upload(ctx, res.Msg.GetSignedUploadUrl(), buf); err != nil {
			return err
		}
		size = int64(buf.Len())
	}

	{
		req := connect.NewRequest(&apiv1.FinalizeArtifactRequest{
			WorkflowRunBackendId:    ids.workflowRunBackendId,
			WorkflowJobRunBackendId: ids.workflowJobRunBackendId,
			Name:                    name,
			Size:                    size,
		})

		res, err := apic.FinalizeArtifact(ctx, req)
		if err != nil {
			return err
		}
		if !res.Msg.GetOk() {
			return errors.New("response is not ok")
		}
	}

	return nil
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
