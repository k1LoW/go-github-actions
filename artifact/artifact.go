package artifact

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const apiVersion = "6.0-preview"

func Upload(ctx context.Context, name string, files []string) error {
	c, err := createContainerForArtifact(ctx, name)
	if err != nil {
		return err
	}
	if err := uploadFiles(ctx, c.FileContainerResourceURL, files); err != nil {
		return err
	}
	return nil
}

type ContainerResponce struct {
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

func createContainerForArtifact(ctx context.Context, name string) (*ContainerResponce, error) {
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
		"POST",
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
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := &ContainerResponce{}
	if err := json.Unmarshal(body, res); err != nil {
		return nil, err
	}

	return res, nil
}

func uploadFiles(ctx context.Context, uploadURL string, files []string) error {
	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			return err
		}
		defer file.Close()
		u := fmt.Sprintf("%s?itemPath=%s", uploadURL, f)
		req, err := http.NewRequest(
			"POST",
			u,
			file,
		)
		req.Header.Set("Content-Type", "application/octet-stream")
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
	}
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
