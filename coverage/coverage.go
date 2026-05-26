// Package coverage uploads a Cobertura XML coverage report to GitHub's code
// coverage API. It is a Go port of github.com/actions/upload-code-coverage.
package coverage

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// ErrSkipped is returned when the upload is intentionally skipped (e.g. for
// merge_group events or pull requests from forks).
var ErrSkipped = errors.New("coverage upload skipped")

// Upload reads the file at path and uploads it as a Cobertura XML coverage
// report. language is a Linguist language name (e.g. "Go"), label is a free
// form label such as "code-coverage/go".
func Upload(ctx context.Context, file, language, label string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return UploadReader(ctx, f, language, label)
}

// UploadReader uploads the coverage report read from content.
func UploadReader(ctx context.Context, content io.Reader, language, label string) error {
	if err := checkSkip(); err != nil {
		return err
	}

	commitOID, ref, prNumber, err := resolveRevision(ctx)
	if err != nil {
		return err
	}

	encoded, err := encodeReport(content)
	if err != nil {
		return err
	}

	payload := map[string]any{
		"commit_oid":      commitOID,
		"coverage_report": encoded,
		"language_name":   language,
		"label":           label,
	}
	switch {
	case prNumber > 0:
		payload["pull_request_number"] = prNumber
	case ref != "":
		payload["ref"] = ref
	default:
		return errors.New("either pull request number or ref must be resolvable from GITHUB_* environment")
	}

	return upload(ctx, payload)
}

func encodeReport(r io.Reader) (string, error) {
	var gz bytes.Buffer
	zw := gzip.NewWriter(&gz)
	if _, err := io.Copy(zw, r); err != nil {
		return "", err
	}
	if err := zw.Close(); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(gz.Bytes()), nil
}

func upload(ctx context.Context, payload map[string]any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	repo := os.Getenv("GITHUB_REPOSITORY")
	if repo == "" {
		return errors.New("env GITHUB_REPOSITORY is not set")
	}
	apiURL := os.Getenv("GITHUB_API_URL")
	if apiURL == "" {
		apiURL = "https://api.github.com"
	}
	token := firstNonEmpty(os.Getenv("GH_TOKEN"), os.Getenv("GITHUB_TOKEN"))
	if token == "" {
		return errors.New("env GH_TOKEN or GITHUB_TOKEN must be set")
	}

	url := fmt.Sprintf("%s/repos/%s/code-coverage/report", strings.TrimRight(apiURL, "/"), repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("coverage upload request failed: %w", err)
	}
	defer res.Body.Close()
	respBody, _ := io.ReadAll(res.Body)

	switch {
	case res.StatusCode == http.StatusCreated:
		return nil
	case res.StatusCode == http.StatusOK:
		// Accepted but not stored (e.g. commit is not the latest on the branch).
		return fmt.Errorf("coverage upload returned HTTP 200 (report not stored): %s", extractMessage(respBody))
	case res.StatusCode == http.StatusForbidden && strings.Contains(strings.ToLower(string(respBody)), "not authorized"):
		return fmt.Errorf("coverage upload returned HTTP %d. Ensure the calling job has 'code-quality: write' permission", res.StatusCode)
	default:
		return fmt.Errorf("coverage upload failed (HTTP %d): %s", res.StatusCode, extractMessage(respBody))
	}
}

func extractMessage(body []byte) string {
	var v struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &v); err == nil && v.Message != "" {
		return v.Message
	}
	return string(body)
}

func firstNonEmpty(vs ...string) string {
	for _, v := range vs {
		if v != "" {
			return v
		}
	}
	return ""
}
