package attest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	protobundle "github.com/sigstore/protobuf-specs/gen/pb-go/bundle/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

func (a *Attester) store(ctx context.Context, b *protobundle.Bundle) error {
	repo := os.Getenv("GITHUB_REPOSITORY")
	if repo == "" {
		return errors.New("env GITHUB_REPOSITORY is only available from the context of an action")
	}
	apiURL := os.Getenv("GITHUB_API_URL")
	if apiURL == "" {
		apiURL = "https://api.github.com"
	}
	bundleJSON, err := protojson.Marshal(b)
	if err != nil {
		return err
	}
	body, err := json.Marshal(map[string]json.RawMessage{"bundle": bundleJSON})
	if err != nil {
		return err
	}
	u := fmt.Sprintf("%s/repos/%s/attestations", strings.TrimSuffix(apiURL, "/"), repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body)) //nolint:gosec // The URL is provided by the runner.
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+a.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	res, err := a.httpClient.Do(req) //nolint:gosec // The URL is provided by the runner.
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		msg, _ := io.ReadAll(io.LimitReader(res.Body, 1024))
		return fmt.Errorf("failed to store attestation: %s: %s", res.Status, bytes.TrimSpace(msg))
	}
	return nil
}
