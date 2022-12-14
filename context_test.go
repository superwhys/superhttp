package superhttp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
)

func TestClientContext(t *testing.T) {
	t.Run("testContextFetch", func(t *testing.T) {
		resp, err := Cli.Start().
			Method("GET").
			URL(fmt.Sprintf("%v/%v", srvApi, "test_get")).
			Headers(DefaultJsonHeader()).
			Fetch(context.Background()).
			BodyString()
		assert.Nil(t, err)
		assert.Equal(t, `{"message":"do success"}`, resp)
	})
}
