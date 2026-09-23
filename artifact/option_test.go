package artifact

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
)

type fakeAttester struct {
	err    error
	called bool
	name   string
	digest map[string]string
}

func (f *fakeAttester) Attest(_ context.Context, name string, digest map[string]string) error {
	f.called = true
	f.name = name
	f.digest = digest
	return f.err
}

func TestAttest(t *testing.T) {
	tests := []struct {
		name        string
		mode        AttestMode
		attester    *fakeAttester
		wantErr     bool
		wantCalled  bool
		wantWarning bool
	}{
		{"none does not attest", AttestNone, &fakeAttester{}, false, false, false},
		{"none ignores attester errors", AttestNone, &fakeAttester{err: errors.New("boom")}, false, false, false},
		{"best effort succeeds", AttestBestEffort, &fakeAttester{}, false, true, false},
		{"best effort warns on failure", AttestBestEffort, &fakeAttester{err: errors.New("boom")}, false, true, true},
		{"best effort warns without attester", AttestBestEffort, nil, false, false, true},
		{"required succeeds", AttestRequired, &fakeAttester{}, false, true, false},
		{"required fails on failure", AttestRequired, &fakeAttester{err: errors.New("boom")}, true, true, false},
		{"required fails without attester", AttestRequired, nil, true, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a Attester
			if tt.attester != nil {
				a = tt.attester
			}
			c := newConfig([]Option{WithAttestation(tt.mode, a)})
			stdout := new(bytes.Buffer)
			c.stdout = stdout
			digest := map[string]string{"sha256": "abc"}

			err := c.attest(context.Background(), "my-artifact", digest)
			if (err != nil) != tt.wantErr {
				t.Errorf("got err %v, want err %v", err, tt.wantErr)
			}
			if tt.attester != nil {
				if tt.attester.called != tt.wantCalled {
					t.Errorf("got called %v, want %v", tt.attester.called, tt.wantCalled)
				}
				if tt.wantCalled && (tt.attester.name != "my-artifact" || tt.attester.digest["sha256"] != "abc") {
					t.Errorf("got subject %s %v", tt.attester.name, tt.attester.digest)
				}
			}
			gotWarning := strings.HasPrefix(stdout.String(), "::warning::")
			if gotWarning != tt.wantWarning {
				t.Errorf("got warning %q, want warning %v", stdout.String(), tt.wantWarning)
			}
		})
	}
}

func TestDefaultConfigDoesNotAttest(t *testing.T) {
	c := newConfig(nil)
	if c.attestMode != AttestNone {
		t.Errorf("got %v, want AttestNone", c.attestMode)
	}
}

func TestEscapeData(t *testing.T) {
	got := escapeData("100%\r\nfailed")
	want := "100%25%0D%0Afailed"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
