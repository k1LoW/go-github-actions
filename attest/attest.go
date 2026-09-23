// Package attest creates GitHub artifact attestations (SLSA build provenance) signed with Sigstore.
package attest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	intoto "github.com/in-toto/attestation/go/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// Attester creates and stores SLSA build provenance attestations on GitHub.
// It satisfies artifact.Attester of github.com/k1LoW/go-github-actions/artifact.
type Attester struct {
	token      string
	httpClient *http.Client
}

// Option is an option for New.
type Option func(*Attester)

// WithToken sets the GitHub token used to store attestations.
// The token requires the attestations:write permission. Defaults to env GITHUB_TOKEN.
func WithToken(token string) Option {
	return func(a *Attester) {
		a.token = token
	}
}

// WithHTTPClient sets the HTTP client used for GitHub API requests.
func WithHTTPClient(c *http.Client) Option {
	return func(a *Attester) {
		a.httpClient = c
	}
}

// New returns a new Attester.
func New(opts ...Option) *Attester {
	a := &Attester{
		token:      os.Getenv("GITHUB_TOKEN"),
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// Attest creates a SLSA build provenance attestation for the subject, signs it with Sigstore and stores it on GitHub.
func (a *Attester) Attest(ctx context.Context, subjectName string, subjectDigest map[string]string) error {
	if a.token == "" {
		return errors.New("GitHub token is required to store attestations")
	}
	idToken, err := requestIDToken(ctx, a.httpClient, "sigstore")
	if err != nil {
		return err
	}
	// The claims are not verified here because Fulcio verifies the same token before issuing the signing certificate.
	claims, err := parseClaims(idToken)
	if err != nil {
		return err
	}
	predicate, err := buildProvenancePredicate(claims, serverURL())
	if err != nil {
		return err
	}
	statement := &intoto.Statement{
		Type: intoto.StatementTypeUri,
		Subject: []*intoto.ResourceDescriptor{
			{Name: subjectName, Digest: subjectDigest},
		},
		PredicateType: slsaPredicateType,
		Predicate:     predicate,
	}
	if err := statement.Validate(); err != nil {
		return fmt.Errorf("invalid statement: %w", err)
	}
	payload, err := protojson.Marshal(statement)
	if err != nil {
		return err
	}
	b, err := signStatement(ctx, payload, idToken, isPublicRepository())
	if err != nil {
		return fmt.Errorf("failed to sign attestation: %w", err)
	}
	return a.store(ctx, b)
}

func serverURL() string {
	if u := os.Getenv("GITHUB_SERVER_URL"); u != "" {
		return u
	}
	return "https://github.com"
}

