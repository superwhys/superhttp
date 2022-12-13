package superhttp

import (
	"bytes"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

type HandleFunc func(c *Context)
type HandlersChain []HandleFunc

type HandlerGroup struct {
	Handlers HandlersChain
	client   *Client
}

func (group *HandlerGroup) Use(handler ...HandleFunc) {
	group.Handlers = append(group.Handlers, handler...)
}

func HandlerDuration() HandleFunc {
	return func(c *Context) {
		startTime := time.Now()
		c.Next()
		c.duration = time.Now().Sub(startTime)
	}
}

func RequestBodyReaderHandler() HandleFunc {
	return func(c *Context) {
		var bodyReader io.Reader
		if c.body != nil {
			bodyReader = bytes.NewReader(c.body)
		}
		c.bodyReader = bodyReader
	}
}

func RequestParamsHandler() HandleFunc {
	return func(c *Context) {
		if c.params != nil {
			urlParse, err := url.ParseRequestURI(c.url)
			if err != nil {
				c.err = errors.Wrap(err, "parse request url")
			}
			q := urlParse.Query()
			for key, value := range c.params {
				if !q.Has(key) {
					q.Add(key, value)
				}
			}
			urlParse.RawQuery = q.Encode()
			c.url = urlParse.String()
		}
	}
}

func RequestDefaultHeaderHandler() HandleFunc {
	return func(c *Context) {
		if c.header == nil {
			c.header = DefaultJsonHeader()
		}
	}
}

func GenerateRequestHandler() HandleFunc {
	return func(c *Context) {
		req, err := http.NewRequest(c.method, c.url, c.bodyReader)
		if err != nil {
			c.err = errors.Wrap(err, "generate request")
			return
		}
		c.Request = req
		req.Header = c.header.Header
	}
}

func DefaultHTTPHandler() HandleFunc {
	return func(c *Context) {
		resp, err := c.client.Do(c.Request)
		c.Response = resp
		if err != nil {
			c.err = err
		}
	}
}
