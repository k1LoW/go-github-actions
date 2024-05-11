package artifact

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"connectrpc.com/connect"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/k1LoW/go-github-actions/artifact/proto/gen/go/results/api/v1/apiv1connect"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"
)

func newAPIClient() (apiv1connect.ArtifactServiceClient, error) {
	addr := os.Getenv("ACTIONS_RESULTS_URL")
	if addr == "" {
		return nil, errors.New("env ACTIONS_RESULTS_URL is only available from the context of an action")
	}
	apic := apiv1connect.NewArtifactServiceClient(httpClient, fmt.Sprintf("%s%s", addr, "twirp"), connect.WithCodec(&protoJSONCodec{}))
	return apic, nil
}

func upload(ctx context.Context, uploadURL string, content io.Reader) error {
	u, err := url.Parse(uploadURL)
	if err != nil {
		return err
	}
	serviceURL := (&url.URL{
		Scheme:   u.Scheme,
		Host:     u.Host,
		RawQuery: u.Query().Encode(),
	}).String()
	splitted := strings.Split(u.Path, "/")
	containerName := splitted[1]
	blobName := strings.Join(splitted[2:], "/")
	uploadc, err := azblob.NewClientWithNoCredential(serviceURL, &azblob.ClientOptions{})
	if err != nil {
		return err
	}
	if _, err := uploadc.UploadStream(ctx, containerName, blobName, content, &azblob.UploadStreamOptions{}); err != nil {
		return err
	}
	return nil
}

type transport struct {
	t http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("ACTIONS_RUNTIME_TOKEN")))
	return t.t.RoundTrip(req)
}

func newTransport() *transport {
	return &transport{
		t: http.DefaultTransport,
	}
}

var httpClient = &http.Client{
	Transport: newTransport(),
}

var _ connect.Codec = (*protoJSONCodec)(nil)

type protoJSONCodec struct{}

func (c *protoJSONCodec) Name() string { return "json" }

func (c *protoJSONCodec) Marshal(message any) ([]byte, error) {
	protoMessage, ok := message.(proto.Message)
	if !ok {
		return nil, errNotProto(message)
	}
	b, err := protojson.MarshalOptions{UseProtoNames: true, EmitDefaultValues: false}.Marshal(protoMessage)
	return b, err
}

func (c *protoJSONCodec) Unmarshal(binary []byte, message any) error {
	protoMessage, ok := message.(proto.Message)
	if !ok {
		return errNotProto(message)
	}
	if len(binary) == 0 {
		return errors.New("zero-length payload is not a valid JSON object")
	}
	options := protojson.UnmarshalOptions{DiscardUnknown: true}
	err := options.Unmarshal(binary, protoMessage)
	if err != nil {
		return fmt.Errorf("unmarshal into %T: %w", message, err)
	}
	return nil
}

func errNotProto(message any) error {
	if _, ok := message.(protoiface.MessageV1); ok {
		return fmt.Errorf("%T uses github.com/golang/protobuf, but connect-go only supports google.golang.org/protobuf: see https://go.dev/blog/protobuf-apiv2", message)
	}
	return fmt.Errorf("%T doesn't implement proto.Message", message)
}
