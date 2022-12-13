package superhttp

import (
	"golang.org/x/net/context"
	"io"
	"math"
	"net/http"
	"time"
)

const abortIndex = math.MaxInt8 >> 1

type Context struct {
	ctx        context.Context
	conf       *Config
	header     *Header
	method     string
	url        string
	params     Params
	body       []byte
	bodyReader io.Reader

	handlers HandlersChain
	index    int8

	cli      *Client
	client   *http.Client
	Request  *http.Request
	Response *http.Response

	duration time.Duration
	err      error
}

func NewContext(client *Client) *Context {
	return &Context{
		cli:      client,
		conf:     client.conf,
		client:   client.httpClient,
		handlers: client.HandlerGroup.Handlers,
		index:    -1,
	}
}

func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *Context) Abort() {
	c.index = abortIndex
}

func (c *Context) URL(url string) *Context {
	if url != "" {
		c.url = url
	}
	return c
}

func (c *Context) Method(method string) *Context {
	c.method = method
	return c
}

func (c *Context) QueryParams(params Params) *Context {
	if params != nil {
		c.params = params
	}
	return c
}

func (c *Context) Headers(header *Header) *Context {
	if header != nil {
		c.header = header
	}
	return c
}

func (c *Context) Body(body []byte) *Context {
	if body != nil {
		c.body = body
	}
	return c
}

func (c *Context) FormBody(form *Form) *Context {
	if form != nil {
		c.body = []byte(form.Encode())
	}
	return c
}

func (c *Context) Fetch(ctx context.Context) *Response {
	if c.err != nil {
		return &Response{err: c.err}
	}
	c.ctx = ctx
	c.Next()
	return &Response{Response: c.Response, err: c.err}
}
