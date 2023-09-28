package baidu_test

import (
	"context"
	. "llm-clients/baidu"
	"testing"
)

func TestCreateAccessToken(t *testing.T) {
	ctx := context.Background()

	client := NewClient("aiPLbu2x50HyPDGtakmaHDxL", "PqZh8PQCY5OKd767MToE4if5kLkkz2A4", false)
	resp, err := client.CreateAccessToken(ctx)

	if err != nil {
		println(err.Error())
	}

	println(resp.AccessToken)

}
