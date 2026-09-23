package attest

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	protobundle "github.com/sigstore/protobuf-specs/gen/pb-go/bundle/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

func testClaims() *claims {
	return &claims{
		Repository:        "k1LoW/go-github-actions",
		RepositoryID:      "1",
		RepositoryOwnerID: "2",
		Ref:               "refs/heads/main",
		SHA:               "deadbeef",
		EventName:         "push",
		WorkflowRef:       "k1LoW/go-github-actions/.github/workflows/ci.yml@refs/heads/main",
		JobWorkflowRef:    "k1LoW/go-github-actions/.github/workflows/ci.yml@refs/heads/main",
		RunnerEnvironment: "github-hosted",
		RunID:             "100",
		RunAttempt:        "1",
	}
}

func TestBuildProvenancePredicate(t *testing.T) {
	p, err := buildProvenancePredicate(testClaims(), "https://github.com")
	if err != nil {
		t.Fatal(err)
	}
	got, err := protojson.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	want := `{
  "buildDefinition": {
    "buildType": "https://actions.github.io/buildtypes/workflow/v1",
    "externalParameters": {
      "workflow": {
        "ref": "refs/heads/main",
        "repository": "https://github.com/k1LoW/go-github-actions",
        "path": ".github/workflows/ci.yml"
      }
    },
    "internalParameters": {
      "github": {
        "event_name": "push",
        "repository_id": "1",
        "repository_owner_id": "2",
        "runner_environment": "github-hosted"
      }
    },
    "resolvedDependencies": [
      {
        "uri": "git+https://github.com/k1LoW/go-github-actions@refs/heads/main",
        "digest": {"gitCommit": "deadbeef"}
      }
    ]
  },
  "runDetails": {
    "builder": {"id": "https://github.com/k1LoW/go-github-actions/.github/workflows/ci.yml@refs/heads/main"},
    "metadata": {"invocationId": "https://github.com/k1LoW/go-github-actions/actions/runs/100/attempts/1"}
  }
}`
	assertJSONEq(t, got, []byte(want))
}

func TestBuildProvenancePredicateInvalidWorkflowRef(t *testing.T) {
	tests := []struct {
		name        string
		workflowRef string
	}{
		{"no ref", "k1LoW/go-github-actions/.github/workflows/ci.yml"},
		{"other repository", "k1LoW/other/.github/workflows/ci.yml@refs/heads/main"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := testClaims()
			c.WorkflowRef = tt.workflowRef
			if _, err := buildProvenancePredicate(c, "https://github.com"); err == nil {
				t.Error("want error")
			}
		})
	}
}

func TestParseClaims(t *testing.T) {
	payload, err := json.Marshal(testClaims())
	if err != nil {
		t.Fatal(err)
	}
	token := "header." + base64.RawURLEncoding.EncodeToString(payload) + ".signature"
	got, err := parseClaims(token)
	if err != nil {
		t.Fatal(err)
	}
	if *got != *testClaims() {
		t.Errorf("got %+v, want %+v", got, testClaims())
	}
	if _, err := parseClaims("invalid"); err == nil {
		t.Error("want error")
	}
}

func TestRequestIDToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer request-token" {
			t.Errorf("got Authorization %q", got)
		}
		if got := r.URL.Query().Get("audience"); got != "sigstore" {
			t.Errorf("got audience %q", got)
		}
		if got := r.URL.Query().Get("api-version"); got != "2.0" {
			t.Errorf("existing query is lost: %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`{"value":"id-token"}`))
	}))
	defer ts.Close()
	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_URL", ts.URL+"/token?api-version=2.0")
	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN", "request-token")

	got, err := requestIDToken(context.Background(), ts.Client(), "sigstore")
	if err != nil {
		t.Fatal(err)
	}
	if got != "id-token" {
		t.Errorf("got %q, want %q", got, "id-token")
	}
}

func TestRequestIDTokenWithoutPermission(t *testing.T) {
	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_URL", "")
	t.Setenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN", "")
	if _, err := requestIDToken(context.Background(), http.DefaultClient, "sigstore"); err == nil {
		t.Error("want error")
	}
}

func TestStore(t *testing.T) {
	b := &protobundle.Bundle{MediaType: "application/vnd.dev.sigstore.bundle.v0.3+json"}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/repos/k1LoW/go-github-actions/attestations" {
			t.Errorf("got %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer gh-token" {
			t.Errorf("got Authorization %q", got)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		want, err := protojson.Marshal(b)
		if err != nil {
			t.Fatal(err)
		}
		assertJSONEq(t, body, []byte(`{"bundle":`+string(want)+`}`))
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":1}`))
	}))
	defer ts.Close()
	t.Setenv("GITHUB_REPOSITORY", "k1LoW/go-github-actions")
	t.Setenv("GITHUB_API_URL", ts.URL)

	a := New(WithToken("gh-token"), WithHTTPClient(ts.Client()))
	if err := a.store(context.Background(), b); err != nil {
		t.Fatal(err)
	}
}

func TestStoreError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()
	t.Setenv("GITHUB_REPOSITORY", "k1LoW/go-github-actions")
	t.Setenv("GITHUB_API_URL", ts.URL)

	a := New(WithToken("gh-token"), WithHTTPClient(ts.Client()))
	if err := a.store(context.Background(), &protobundle.Bundle{}); err == nil {
		t.Error("want error")
	}
}

func TestAttestWithoutToken(t *testing.T) {
	a := New(WithToken(""))
	if err := a.Attest(context.Background(), "subject", map[string]string{"sha256": "abc"}); err == nil {
		t.Error("want error")
	}
}

func TestIsPublicRepository(t *testing.T) {
	tests := []struct {
		name  string
		event string
		want  bool
	}{
		{"public", `{"repository":{"visibility":"public"}}`, true},
		{"private", `{"repository":{"visibility":"private"}}`, false},
		{"internal", `{"repository":{"visibility":"internal"}}`, false},
		{"unknown", `{}`, false},
		{"broken", `{`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := filepath.Join(t.TempDir(), "event.json")
			if err := os.WriteFile(p, []byte(tt.event), 0o600); err != nil {
				t.Fatal(err)
			}
			t.Setenv("GITHUB_EVENT_PATH", p)
			if got := isPublicRepository(); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
	t.Run("no event", func(t *testing.T) {
		t.Setenv("GITHUB_EVENT_PATH", "")
		if isPublicRepository() {
			t.Error("got true, want false")
		}
	})
}

func TestSigningEndpoints(t *testing.T) {
	tests := []struct {
		name      string
		public    bool
		serverURL string
		want      *endpoints
		wantErr   bool
	}{
		{"public repository", true, "https://github.com", &endpoints{fulcioURL: "https://fulcio.sigstore.dev", rekorURL: "https://rekor.sigstore.dev"}, false},
		{"private repository on github.com", false, "https://github.com", &endpoints{fulcioURL: "https://fulcio.githubapp.com", tsaURL: "https://timestamp.githubapp.com/api/v1/timestamp"}, false},
		{"private repository on ghe.com", false, "https://octo.ghe.com", &endpoints{fulcioURL: "https://fulcio.octo.ghe.com", tsaURL: "https://timestamp.octo.ghe.com/api/v1/timestamp"}, false},
		{"invalid server URL", false, "github.com", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := signingEndpoints(tt.public, tt.serverURL)
			if (err != nil) != tt.wantErr {
				t.Fatalf("got err %v, want err %v", err, tt.wantErr)
			}
			if tt.want != nil && *got != *tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func assertJSONEq(t *testing.T, got, want []byte) {
	t.Helper()
	var g, w any
	if err := json.Unmarshal(got, &g); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(want, &w); err != nil {
		t.Fatal(err)
	}
	gb, _ := json.Marshal(g)
	wb, _ := json.Marshal(w)
	if string(gb) != string(wb) {
		t.Errorf("got %s\nwant %s", gb, wb)
	}
}
