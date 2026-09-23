package attest

import (
	"context"
	"encoding/json"
	"os"

	protobundle "github.com/sigstore/protobuf-specs/gen/pb-go/bundle/v1"
	"github.com/sigstore/sigstore-go/pkg/bundle"
	"github.com/sigstore/sigstore-go/pkg/sign"
)

const (
	publicGoodFulcioURL = "https://fulcio.sigstore.dev"
	publicGoodRekorURL  = "https://rekor.sigstore.dev"
	githubFulcioURL     = "https://fulcio.githubapp.com"
	githubTSAURL        = "https://timestamp.githubapp.com/api/v1/timestamp"
	retries             = 3
)

// signStatement signs the in-toto statement with the same Sigstore instance that @actions/attest selects.
// Public repositories use the Public Good instance and record the signature in Rekor.
// Other repositories use the GitHub instance, which uses a timestamp authority instead of a transparency log
// so that details of private repositories are not published.
func signStatement(ctx context.Context, statement []byte, idToken string, public bool) (*protobundle.Bundle, error) {
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
	if public {
		opts.CertificateProvider = sign.NewFulcio(&sign.FulcioOptions{BaseURL: publicGoodFulcioURL, Retries: retries})
		opts.TransparencyLogs = []sign.Transparency{
			sign.NewRekor(&sign.RekorOptions{BaseURL: publicGoodRekorURL, Retries: retries}),
		}
	} else {
		opts.CertificateProvider = sign.NewFulcio(&sign.FulcioOptions{BaseURL: githubFulcioURL, Retries: retries})
		opts.TimestampAuthorities = []*sign.TimestampAuthority{
			sign.NewTimestampAuthority(&sign.TimestampAuthorityOptions{URL: githubTSAURL, Retries: retries}),
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
