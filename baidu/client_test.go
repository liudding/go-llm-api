package baidu //nolint:testpackage // testing private field

import (
	"context"
	"errors"
	"net/http"
	"testing"
)

var errTestRequestBuilderFailed = errors.New("test request builder failed")

type failingRequestBuilder struct{}

func (*failingRequestBuilder) Build(_ context.Context, _, _ string, _ any, _ http.Header) (*http.Request, error) {
	return nil, errTestRequestBuilderFailed
}

func TestClient(t *testing.T) {

}
