package coverage

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEncodeReport(t *testing.T) {
	src := []byte("<coverage></coverage>")
	got, err := encodeReport(bytes.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	raw, err := base64.StdEncoding.DecodeString(got)
	if err != nil {
		t.Fatalf("base64 decode: %v", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("gzip reader: %v", err)
	}
	defer zr.Close()
	out, err := io.ReadAll(zr)
	if err != nil {
		t.Fatalf("gzip read: %v", err)
	}
	if !bytes.Equal(out, src) {
		t.Errorf("roundtrip mismatch: got %q, want %q", out, src)
	}
}

func TestCheckSkip_MergeGroup(t *testing.T) {
	t.Setenv("GITHUB_EVENT_NAME", "merge_group")
	t.Setenv("GITHUB_EVENT_PATH", "")
	if err := checkSkip(); err == nil || !strings.Contains(err.Error(), "merge_group") {
		t.Errorf("expected merge_group skip, got %v", err)
	}
}

func TestCheckSkip_ForkPR(t *testing.T) {
	dir := t.TempDir()
	ev := map[string]any{
		"pull_request": map[string]any{
			"head": map[string]any{
				"repo": map[string]any{"full_name": "fork/repo"},
			},
		},
	}
	b, _ := json.Marshal(ev)
	path := filepath.Join(dir, "event.json")
	if err := os.WriteFile(path, b, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GITHUB_EVENT_NAME", "pull_request")
	t.Setenv("GITHUB_EVENT_PATH", path)
	t.Setenv("GITHUB_REPOSITORY", "owner/repo")
	if err := checkSkip(); err == nil || !strings.Contains(err.Error(), "fork") {
		t.Errorf("expected fork PR skip, got %v", err)
	}
}

func TestResolveRevision_PullRequest(t *testing.T) {
	dir := t.TempDir()
	ev := map[string]any{
		"pull_request": map[string]any{
			"number": float64(42),
			"head":   map[string]any{"sha": "deadbeef"},
		},
	}
	b, _ := json.Marshal(ev)
	path := filepath.Join(dir, "event.json")
	if err := os.WriteFile(path, b, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GITHUB_EVENT_NAME", "pull_request")
	t.Setenv("GITHUB_EVENT_PATH", path)

	sha, ref, num, err := resolveRevision(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if sha != "deadbeef" || ref != "" || num != 42 {
		t.Errorf("got sha=%q ref=%q num=%d", sha, ref, num)
	}
}

func TestUploadReader_Success(t *testing.T) {
	var received map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s", r.Method)
		}
		if r.URL.Path != "/repos/owner/repo/code-coverage/report" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("auth header = %s", got)
		}
		_ = json.NewDecoder(r.Body).Decode(&received)
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	dir := t.TempDir()
	ev := map[string]any{
		"pull_request": map[string]any{
			"number": float64(7),
			"head":   map[string]any{"sha": "abc123"},
		},
	}
	b, _ := json.Marshal(ev)
	path := filepath.Join(dir, "event.json")
	if err := os.WriteFile(path, b, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GITHUB_EVENT_NAME", "pull_request")
	t.Setenv("GITHUB_EVENT_PATH", path)
	t.Setenv("GITHUB_REPOSITORY", "owner/repo")
	t.Setenv("GITHUB_API_URL", srv.URL)
	t.Setenv("GH_TOKEN", "test-token")
	t.Setenv("GITHUB_TOKEN", "")

	err := UploadReader(context.Background(), strings.NewReader("<coverage/>"), "Go", "code-coverage/go")
	if err != nil {
		t.Fatalf("UploadReader: %v", err)
	}
	if received["language_name"] != "Go" || received["label"] != "code-coverage/go" {
		t.Errorf("payload language/label wrong: %v", received)
	}
	if received["commit_oid"] != "abc123" {
		t.Errorf("commit_oid = %v", received["commit_oid"])
	}
	if pr, _ := received["pull_request_number"].(float64); pr != 7 {
		t.Errorf("pull_request_number = %v", received["pull_request_number"])
	}
	if _, ok := received["coverage_report"].(string); !ok {
		t.Errorf("coverage_report missing")
	}
}

func TestUploadReader_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"message":"invalid report"}`))
	}))
	defer srv.Close()

	dir := t.TempDir()
	ev := map[string]any{
		"pull_request": map[string]any{
			"number": float64(1),
			"head":   map[string]any{"sha": "sha"},
		},
	}
	b, _ := json.Marshal(ev)
	path := filepath.Join(dir, "event.json")
	if err := os.WriteFile(path, b, 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GITHUB_EVENT_NAME", "pull_request")
	t.Setenv("GITHUB_EVENT_PATH", path)
	t.Setenv("GITHUB_REPOSITORY", "owner/repo")
	t.Setenv("GITHUB_API_URL", srv.URL)
	t.Setenv("GH_TOKEN", "test-token")

	err := UploadReader(context.Background(), strings.NewReader("<coverage/>"), "Go", "code-coverage/go")
	if err == nil || !strings.Contains(err.Error(), "invalid report") {
		t.Errorf("expected invalid report error, got %v", err)
	}
}

func TestUpload_FileNotFound(t *testing.T) {
	t.Setenv("GITHUB_EVENT_NAME", "")
	t.Setenv("GITHUB_EVENT_PATH", "")
	err := Upload(context.Background(), "/no/such/file.xml", "Go", "code-coverage/go")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestUploadIntegration(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "" {
		t.Skip("Not running on GitHub Actions")
	}
	if err := Upload(context.Background(), "testdata/cobertura.xml", "Go", "code-coverage/go-test"); err != nil {
		t.Error(err)
	}
}
