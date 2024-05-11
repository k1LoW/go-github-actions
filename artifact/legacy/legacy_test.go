package legacy

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestGetArtifactURL(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "" {
		t.Skip("Not running on GitHub Actions")
	}
	u, err := getArtifactURL()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(u, "/_apis/pipelines") {
		t.Errorf("invalid URL. got %s", u)
	}
}

func TestCreateContainerForArtifact(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "" {
		t.Skip("Not running on GitHub Actions")
	}
	got, err := createContainerForArtifact(context.TODO(), "TestCreateContainerForArtifact")
	if err != nil {
		t.Error(err)
	}
	if want := "actions_storage"; got.Type != want {
		t.Errorf("got %v\nwant %v", got.Name, want)
	}
	if want := "TestCreateContainerForArtifact"; got.Name != want {
		t.Errorf("got %v\nwant %v", got.Name, want)
	}
	if want := ".actions.githubusercontent.com/"; !strings.Contains(got.FileContainerResourceURL, want) {
		t.Errorf("got %v\nwant %v*", got.FileContainerResourceURL, want)
	}
}

func TestUpload(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "" {
		t.Skip("Not running on GitHub Actions")
	}
	if err := Upload(context.TODO(), "TestUploadLegacy", "artifact/testdata/test.txt", strings.NewReader("hello artifact 3\n")); err != nil {
		t.Error(err)
	}
}

func TestUploadLargeContent(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "" {
		t.Skip("Not running on GitHub Actions")
	}
	const (
		owner = "k1LoW"
		repo  = "go-github-actions"
	)
	ctx := context.TODO()
	s := strings.Repeat("0123456789\n", 1024*1024*10)
	name := "TestUploadLargeContentLegacy"
	if err := Upload(ctx, name, "artifact/testdata/large.txt", strings.NewReader(s)); err != nil {
		t.Error(err)
	}
}

func TestUploadFiles(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") == "" {
		t.Skip("Not running on GitHub Actions")
	}
	files := []string{
		"../testdata/test2.txt",
		"../testdata/test3.txt",
	}
	if err := UploadFiles(context.TODO(), "TestUploadFilesLegacy", files); err != nil {
		t.Error(err)
	}
}
