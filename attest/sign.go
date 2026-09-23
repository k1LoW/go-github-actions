package attest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	protobundle "github.com/sigstore/protobuf-specs/gen/pb-go/bundle/v1"
	"github.com/sigstore/sigstore-go/pkg/bundle"
	"github.com/sigstore/sigstore-go/pkg/sign"
)

const (
	publicGoodFulcioURL = "https://fulcio.sigstore.dev"
	publicGoodRekorURL  = "https://rekor.sigstore.dev"
	retries             = 3
)

type endpoints struct {
	fulcioURL string
	rekorURL  string
	tsaURL    string
}

// signingEndpoints selects the same Sigstore instance as @actions/attest.
// Public repositories use the Public Good instance and record the signature in Rekor.
// Other repositories use the GitHub instance of the server, which uses a timestamp authority instead of a
// transparency log so that details of private repositories are not published.
func signingEndpoints(public bool, serverURL string) (*endpoints, error) {
	if public {
		return &endpoints{fulcioURL: publicGoodFulcioURL, rekorURL: publicGoodRekorURL}, nil
	}
	u, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}
	host := u.Hostname()
	if host == "" {
		return nil, fmt.Errorf("invalid server URL: %s", serverURL)
	}
	if host == "github.com" {
		host = "githubapp.com"
	}
	return &endpoints{
		fulcioURL: "https://fulcio." + host,
		tsaURL:    "https://timestamp." + host + "/api/v1/timestamp",
	}, nil
}

func signStatement(ctx context.Context, statement []byte, idToken string, eps *endpoints) (*protobundle.Bundle, error) {
	keypair, err := sign.NewEphemeralKeypair(nil)
	if err != nil {
		return nil, err
	}
	content := &sign.DSSEData{
		Data:        statement,
		PayloadType: bundle.IntotoMediaType,
	}
	opts := sign.BundleOptions{
		Context:                    ctx,
		CertificateProviderOptions: &sign.CertificateProviderOptions{IDToken: idToken},
	}
	opts.CertificateProvider = sign.NewFulcio(&sign.FulcioOptions{BaseURL: eps.fulcioURL, Retries: retries})
	if eps.rekorURL != "" {
		opts.TransparencyLogs = []sign.Transparency{
			sign.NewRekor(&sign.RekorOptions{BaseURL: eps.rekorURL, Retries: retries}),
		}
	}
	if eps.tsaURL != "" {
		opts.TimestampAuthorities = []*sign.TimestampAuthority{
			sign.NewTimestampAuthority(&sign.TimestampAuthorityOptions{URL: eps.tsaURL, Retries: retries}),
		}
	}
	return sign.Bundle(content, keypair, opts)
}

// isPublicRepository reports whether the repository is public by reading the event payload.
// It falls back to false so that a repository of unknown visibility is never recorded in the public transparency log.
func isPublicRepository() bool {
	p := os.Getenv("GITHUB_EVENT_PATH")
	if p == "" {
		return false
	}
	b, err := os.ReadFile(p) //nolint:gosec // The path is provided by the runner.
	if err != nil {
		return false
	}
	var event struct {
		Repository struct {
			Visibility string `json:"visibility"`
		} `json:"repository"`
	}
	if err := json.Unmarshal(b, &event); err != nil {
		return false
	}
	return event.Repository.Visibility == "public"
}
