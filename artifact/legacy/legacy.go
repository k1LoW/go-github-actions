package legacy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const apiVersion = "6.0-preview"
const uploadChunkSize = 8 * 1024 * 1024 // 8 MB

// Upload content as GitHub Actions artifact.
func Upload(ctx context.Context, name, fp string, content io.Reader) error {
	c, err := createContainerForArtifact(ctx, name)
	if err != nil {
		return err
	}

	size, err := upload(ctx, name, c.FileContainerResourceURL, fp, content)
	if err != nil {
		return err
	}

	if err := patchArtifactSize(ctx, name, size); err != nil {
		return err
	}

	return nil
}

// UploadFiles as GitHub Actions artifact.
func UploadFiles(ctx context.Context, name string, files []string) error {
	c, err := createContainerForArtifact(ctx, name)
	if err != nil {
		return err
	}

	total, err := uploadFiles(ctx, name, c.FileContainerResourceURL, files)
	if err != nil {
		return err
	}

	if err := patchArtifactSize(ctx, name, total); err != nil {
		return err
	}

	return nil
}

type containerResponce struct {
	ContainerID              int         `json:"containerId"`
	Size                     int         `json:"size"`
	SignedContent            interface{} `json:"signedContent"`
	FileContainerResourceURL string      `json:"fileContainerResourceUrl"`
	Type                     string      `json:"type"`
	Name                     string      `json:"name"`
	URL                      string      `json:"url"`
	ExpiresOn                time.Time   `json:"expiresOn"`
	Items                    interface{} `json:"items"`
}

func createContainerForArtifact(ctx context.Context, name string) (*containerResponce, error) {
	param := map[string]string{
		"Type": "actions_storage",
		"Name": name,
	}

	u, err := getArtifactURL()
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(&param)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		u,
		bytes.NewReader(b),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", fmt.Sprintf("application/json;api-version=%s", apiVersion))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ACTIONS_RUNTIME_TOKEN")))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := &containerResponce{}
	if err := json.Unmarshal(body, res); err != nil {
		return nil, err
	}

	return res, nil
}

func upload(ctx context.Context, name, ep, fp string, content io.Reader) (int, error) {
	u, err := url.Parse(ep)
	if err != nil {
		return 0, err
	}
	q := u.Query()
	q.Set("itemPath", filepath.Join(name, fp))
	q.Set("api-version", apiVersion)
	u.RawQuery = q.Encode()
	body := &bytes.Buffer{}
	if _, err = io.Copy(body, content); err != nil {
		return 0, err
	}
	max := body.Len()
	buf := make([]byte, 0, uploadChunkSize)
	start := 0
	client := &http.Client{}
	for {
		n, err := body.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if errors.Is(err, io.EOF) {
				break
			}
			return 0, err
		}
		end := start + n - 1
		req, err := createRequest(u, start, end, max, bytes.NewReader(buf))
		if err != nil {
			return 0, err
		}
		resp, err := client.Do(req)
		if err != nil {
			return 0, err
		}
		if _, err := io.ReadAll(resp.Body); err != nil {
			return 0, err
		}
		if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
			return 0, errors.New(resp.Status)
		}
		start = start + n
		if err != nil && !errors.Is(err, io.EOF) {
			return 0, err
		}
	}

	return max, nil
}

func patchArtifactSize(ctx context.Context, name string, size int) error {
	e, err := getArtifactURL()
	if err != nil {
		return err
	}
	u, err := url.Parse(e)
	if err != nil {
		return err
	}
	q := u.Query()
	q.Set("artifactName", name)
	q.Set("api-version", apiVersion)
	u.RawQuery = q.Encode()

	param := map[string]int{
		"Size": size,
	}
	b, err := json.Marshal(&param)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(
		http.MethodPatch,
		u.String(),
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", fmt.Sprintf("application/json;api-version=%s", apiVersion))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ACTIONS_RUNTIME_TOKEN")))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if _, err := io.ReadAll(resp.Body); err != nil {
		return err
	}
	return nil
}

func uploadFiles(ctx context.Context, name, ep string, files []string) (int, error) {
	total := 0
	for _, f := range files {
		a, err := filepath.Abs(f)
		if err != nil {
			return 0, err
		}

		rel, err := filepath.Rel(os.Getenv("GITHUB_WORKSPACE"), a)
		if err != nil {
			return 0, err
		}

		file, err := os.Open(filepath.Clean(f))
		if err != nil {
			return 0, err
		}
		size, err := upload(ctx, name, ep, rel, file)
		if err != nil {
			_ = file.Close()
			return 0, err
		}
		total += size
		if err := file.Close(); err != nil {
			return 0, err
		}

	}
	return total, nil
}

func createRequest(u *url.URL, start, end, max int, b io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(
		http.MethodPut,
		u.String(),
		b,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", fmt.Sprintf("application/json;api-version=%s", apiVersion))
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ACTIONS_RUNTIME_TOKEN")))
	req.Header.Set("Content-Length", fmt.Sprintf("%d", end-start+1))
	req.Header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, max))

	return req, nil
}

func getArtifactURL() (string, error) {
	if os.Getenv("ACTIONS_RUNTIME_URL") == "" {
		return "", errors.New("env ACTIONS_RUNTIME_URL is only available from the context of an action")
	}
	if os.Getenv("GITHUB_RUN_ID") == "" {
		return "", errors.New("env GITHUB_RUN_ID is only available from the context of an action")
	}
	return fmt.Sprintf("%s_apis/pipelines/workflows/%s/artifacts?api-version=%s", os.Getenv("ACTIONS_RUNTIME_URL"), os.Getenv("GITHUB_RUN_ID"), apiVersion), nil
}
