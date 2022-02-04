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
	createContainerForArtifact(context.TODO(), "TestCreateContainerForArtifact")

	t.Fatal(1)
}
