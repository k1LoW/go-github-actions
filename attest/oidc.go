package attest

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type claims struct {
	Repository        string `json:"repository"`
	RepositoryID      string `json:"repository_id"`
	RepositoryOwnerID string `json:"repository_owner_id"`
	Ref               string `json:"ref"`
	SHA               string `json:"sha"`
	EventName         string `json:"event_name"`
	WorkflowRef       string `json:"workflow_ref"`
	JobWorkflowRef    string `json:"job_workflow_ref"`
	RunnerEnvironment string `json:"runner_environment"`
	RunID             string `json:"run_id"`
	RunAttempt        string `json:"run_attempt"`
}

func requestIDToken(ctx context.Context, c *http.Client, audience string) (string, error) {
	reqURL := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL")
	reqToken := os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")
	if reqURL == "" || reqToken == "" {
		return "", errors.New("env ACTIONS_ID_TOKEN_REQUEST_URL and ACTIONS_ID_TOKEN_REQUEST_TOKEN are required. Grant the id-token:write permission to the workflow")
	}
	u, err := url.Parse(reqURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("audience", audience)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil) //nolint:gosec // The URL is provided by the runner.
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+reqToken)
	req.Header.Set("Accept", "application/json")
	res, err := c.Do(req) //nolint:gosec // The URL is provided by the runner.
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to request ID token: %s", res.Status)
	}
	var body struct {
		Value string `json:"value"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return "", err
	}
	if body.Value == "" {
		return "", errors.New("empty ID token")
	}
	return body.Value, nil
}

func parseClaims(token string) (*claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid ID token")
	}
	b, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid ID token: %w", err)
	}
	c := &claims{}
	if err := json.Unmarshal(b, c); err != nil {
		return nil, fmt.Errorf("invalid ID token: %w", err)
	}
	return c, nil
}
