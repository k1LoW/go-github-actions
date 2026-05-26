package coverage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// checkSkip returns ErrSkipped when the current event should not produce an
// upload (merge queue runs, and pull requests from forked repositories).
func checkSkip() error {
	if os.Getenv("GITHUB_EVENT_NAME") == "merge_group" {
		return fmt.Errorf("%w: merge_group event", ErrSkipped)
	}
	ev, err := readEvent()
	if err != nil {
		return err
	}
	headRepo := pullRequestHeadRepoFullName(ev)
	if headRepo != "" && headRepo != os.Getenv("GITHUB_REPOSITORY") {
		return fmt.Errorf("%w: pull request from fork %s", ErrSkipped, headRepo)
	}
	return nil
}

// resolveRevision determines the commit OID, ref and pull request number to
// report based on the active GitHub Actions event.
func resolveRevision(ctx context.Context) (commitOID, ref string, prNumber int, err error) {
	eventName := os.Getenv("GITHUB_EVENT_NAME")
	ev, err := readEvent()
	if err != nil {
		return "", "", 0, err
	}

	if eventName == "pull_request" || eventName == "pull_request_target" {
		pr, ok := ev["pull_request"].(map[string]any)
		if !ok {
			return "", "", 0, fmt.Errorf("event payload has no pull_request object")
		}
		head, _ := pr["head"].(map[string]any)
		sha, _ := head["sha"].(string)
		num, _ := pr["number"].(float64)
		return sha, "", int(num), nil
	}

	commitOID = os.Getenv("GITHUB_SHA")
	ref = os.Getenv("GITHUB_REF")
	prNumber, _ = lookupOpenPullRequest(ctx, os.Getenv("GITHUB_REPOSITORY"), os.Getenv("GITHUB_REF_NAME"))
	return commitOID, ref, prNumber, nil
}

func readEvent() (map[string]any, error) {
	path := os.Getenv("GITHUB_EVENT_PATH")
	if path == "" {
		return map[string]any{}, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]any{}, nil
		}
		return nil, err
	}
	if len(b) == 0 {
		return map[string]any{}, nil
	}
	var ev map[string]any
	if err := json.Unmarshal(b, &ev); err != nil {
		return nil, fmt.Errorf("parse GITHUB_EVENT_PATH: %w", err)
	}
	return ev, nil
}

func pullRequestHeadRepoFullName(ev map[string]any) string {
	pr, ok := ev["pull_request"].(map[string]any)
	if !ok {
		return ""
	}
	head, ok := pr["head"].(map[string]any)
	if !ok {
		return ""
	}
	repo, ok := head["repo"].(map[string]any)
	if !ok {
		return ""
	}
	name, _ := repo["full_name"].(string)
	return name
}

// lookupOpenPullRequest finds the open pull request whose head matches branch.
// It returns 0 (no error) when none is found so push-event uploads still
// proceed with ref-based association.
func lookupOpenPullRequest(ctx context.Context, repo, branch string) (int, error) {
	if repo == "" || branch == "" {
		return 0, nil
	}
	token := firstNonEmpty(os.Getenv("GH_TOKEN"), os.Getenv("GITHUB_TOKEN"))
	if token == "" {
		return 0, nil
	}
	apiURL := os.Getenv("GITHUB_API_URL")
	if apiURL == "" {
		apiURL = "https://api.github.com"
	}
	owner := strings.SplitN(repo, "/", 2)[0]
	q := url.Values{}
	q.Set("state", "open")
	q.Set("head", fmt.Sprintf("%s:%s", owner, branch))
	u := fmt.Sprintf("%s/repos/%s/pulls?%s", strings.TrimRight(apiURL, "/"), repo, q.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return 0, nil
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return 0, nil
	}
	var prs []struct {
		Number int `json:"number"`
	}
	if err := json.NewDecoder(res.Body).Decode(&prs); err != nil {
		return 0, nil
	}
	if len(prs) == 0 {
		return 0, nil
	}
	return prs[0].Number, nil
}
