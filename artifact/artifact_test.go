package artifact

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestGetArtifactURL(t *testing.T) {
	if os.Getenv("CI") == "" {
		t.Skip("env CI is not set")
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
	if os.Getenv("CI") == "" {
		t.Skip("env CI is not set")
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
	if want := "https://pipelines.actions.githubusercontent.com/"; !strings.HasPrefix(got.FileContainerResourceURL, want) {
		t.Errorf("got %v\nwant %v*", got.FileContainerResourceURL, want)
	}
}

func TestUpload(t *testing.T) {
	if os.Getenv("CI") == "" {
		t.Skip("env CI is not set")
	}
	files := []string{"artifact_test.go"}
	if err := Upload(context.TODO(), "TestUpload", files); err != nil {
		t.Error(err)
	}
}
