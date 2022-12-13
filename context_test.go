package superhttp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
)

var (
	srvApi = "http://127.0.0.1:29940"
)

func TestClientContext(t *testing.T) {
	t.Run("testContextFetch", func(t *testing.T) {
		resp, err := Cli.Start().
			Method("GET").
			URL(fmt.Sprintf("%v/%v", srvApi, "test_get")).
			Headers(DefaultJsonHeader()).
			Fetch(context.Background()).
			Body()
		assert.Nil(t, err)
		assert.Equal(t, `{"message":"do success"}`, resp)
	})
}
