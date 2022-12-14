package superhttp

import (
	"crypto/tls"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	HandlerGroup
	conf       *Config
	httpClient *http.Client
	isDefault  bool
}

func newClient(conf *Config) (*Client, error) {
	var (
		transportProxy  func(*http.Request) (*url.URL, error)
		tlsClientConfig *tls.Config
	)

	if conf.proxy != "" {
		transportProxy = func(_ *http.Request) (*url.URL, error) {
			return url.Parse(conf.proxy)
		}
		tlsClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	} else {
		transportProxy = http.ProxyFromEnvironment
	}
	cli := &http.Client{
		Transport: &http.Transport{
			Proxy: transportProxy,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       tlsClientConfig,
		},
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       conf.requestTimeOut,
	}
	return &Client{
		conf:       conf,
		httpClient: cli,
	}, nil
}

func New(conf *Config) *Client {
	if conf == nil {
		panic("need httpclient config")
	}
	if conf.requestTimeOut == 0 {
		panic("httpclient request timeout is not legal")
	}
	client, err := newClient(conf)
	if err != nil {
		panic(errors.Wrap(err, "new httpclient failed"))
	}
	client.isDefault = false
	return client
}

func Default() *Client {
	conf := &Config{
		requestTimeOut: 10 * time.Second,
	}

	cli := New(conf)
	cli.Use(
		HandlerDuration(),
		RequestDefaultHeaderHandler(),
		RequestParamsHandler(),
		RequestBodyReaderHandler(),
		DefaultHTTPHandler(),
	)
	cli.isDefault = true
	return cli
}

func (cli *Client) Use(handler ...HandleFunc) {
	cli.HandlerGroup.Use(handler...)
}

func (cli *Client) Start() *Context {
	ctx := NewContext(cli)
	return ctx
}

func (cli *Client) DoRequest(ctx context.Context, url, method string, queryParams Params, header *Header, body []byte, callBack ...HandleFunc) *Response {
	cli.Use(callBack...)
	return cli.
		Start().
		Method(method).
		URL(url).
		QueryParams(queryParams).
		Headers(header).
		Body(body).
		Fetch(ctx)
}

func (cli *Client) Get(ctx context.Context, url string, queryParams Params, header *Header, callBack ...HandleFunc) *Response {
	cli.Use(callBack...)
	return cli.DoRequest(ctx, url, http.MethodGet, queryParams, header, nil, callBack...)
}

func (cli *Client) Post(ctx context.Context, url string, queryParams Params, header *Header, body []byte, callBack ...HandleFunc) *Response {
	cli.Use(callBack...)
	return cli.DoRequest(ctx, url, http.MethodPost, queryParams, header, body, callBack...)
}
