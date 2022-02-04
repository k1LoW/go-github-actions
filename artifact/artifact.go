package artifact

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const apiVersion = "6.0-preview"

func Upload(ctx context.Context, name string, r io.Reader) error {
	err := createContainerForArtifact(ctx, name)
	if err != nil {
		return err
	}
	return nil
}

func createContainerForArtifact(ctx context.Context, name string) error {
	values := url.Values{}
	values.Set("Type", "actions_storage")
	values.Add("Name", name)

	u, err := getArtifactURL()
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		u,
		strings.NewReader(values.Encode()),
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
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(body))

	return nil
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
