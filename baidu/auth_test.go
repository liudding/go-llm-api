package baidu_test

import (
	"context"
	. "github.com/liudding/go-llm-api/baidu"
	"testing"
)

func TestCreateAccessToken(t *testing.T) {
	ctx := context.Background()

	client := NewClient("xxxx", "yyyy", false)
	resp, err := client.CreateAccessToken(ctx)

	if err != nil {
		println(err.Error())
	}

	println(resp.AccessToken)

}
