package artifact

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// AttestMode controls whether an artifact attestation is created on upload.
type AttestMode int

const (
	// AttestNone does not create an attestation.
	AttestNone AttestMode = iota
	// AttestBestEffort tries to create an attestation and only emits a warning if it fails.
	AttestBestEffort
	// AttestRequired creates an attestation and fails the upload if it fails.
	AttestRequired
)

// Attester creates an attestation for an artifact.
// github.com/k1LoW/go-github-actions/attest provides an implementation.
//
// The signature uses only built-in types so that implementations do not need to import this package.
type Attester interface {
	Attest(ctx context.Context, subjectName string, subjectDigest map[string]string) error
}

// Option is an option for Upload and UploadFiles.
type Option func(*config)

// WithAttestation sets the attestation mode and the Attester used to create attestations.
func WithAttestation(mode AttestMode, attester Attester) Option {
	return func(c *config) {
		c.attestMode = mode
		c.attester = attester
	}
}

type config struct {
	attestMode AttestMode
	attester   Attester
	stdout     io.Writer
}

func newConfig(opts []Option) *config {
	c := &config{
		attestMode: AttestNone,
		stdout:     os.Stdout,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *config) attest(ctx context.Context, name string, digest map[string]string) error {
	switch c.attestMode {
	case AttestNone:
		return nil
	case AttestBestEffort, AttestRequired:
	default:
		return fmt.Errorf("unknown attest mode: %d", c.attestMode)
	}
	if c.attester == nil {
		return c.handleAttestError(errors.New("no attester is set"))
	}
	return c.handleAttestError(c.attester.Attest(ctx, name, digest))
}

func (c *config) handleAttestError(err error) error {
	if err == nil {
		return nil
	}
	switch c.attestMode {
	case AttestRequired:
		return fmt.Errorf("failed to attest artifact: %w", err)
	case AttestBestEffort:
		// Returning nil silently would hide the failure, so surface it as a workflow warning annotation.
		_, _ = fmt.Fprintf(c.stdout, "::warning::failed to attest artifact: %s\n", escapeData(err.Error()))
	case AttestNone:
	default:
		return fmt.Errorf("unknown attest mode %d: %w", c.attestMode, err)
	}
	return nil
}

func escapeData(s string) string {
	return strings.NewReplacer("%", "%25", "\r", "%0D", "\n", "%0A").Replace(s)
}
