package superhttp

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/url"
	"os"
	"testing"
	"time"
)

var (
	Cli    *Client
	srvApi = "http://127.0.0.1:29940"
)

type params struct {
	Name string `json:"name" form:"name"`
}

func TestMain(m *testing.M) {
	Cli = Default()

	go func() {
		r := gin.Default()
		r.GET("/test_get", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "do success",
			})
		})
		r.GET("/test_get_params", func(c *gin.Context) {
			p := &params{}
			err := c.ShouldBind(p)
			if err != nil {
				c.JSON(400, gin.H{"err": err.Error()})
				return
			}
			fmt.Printf("params:%v\n", p)
			c.JSON(200, gin.H{"message": fmt.Sprintf("hello %v", p.Name)})
		})
		panic(r.Run(":29940"))
	}()
	time.Sleep(2 * time.Second)

	os.Exit(m.Run())
}

func addParamsHandler(c *Context) {
	fmt.Printf("url:%v\n", c.url)
	urlParse, err := url.ParseRequestURI(c.url)
	if err != nil {
		c.err = errors.Wrap(err, "parse request url")
		return
	}
	q := urlParse.Query()
	p := Params{"name": "superwhys"}
	for key, value := range p {
		if !q.Has(key) {
			q.Add(key, value)
		}
	}
	urlParse.RawQuery = q.Encode()
	c.url = urlParse.String()
}

func TestNewClientGet(t *testing.T) {
	t.Run("testNewClientGet", func(t *testing.T) {
		newCli := New(&Config{
			requestTimeOut: 5 * time.Second,
		})
		newCli.Use(addParamsHandler, DefaultHTTPHandler())
		resp, err := newCli.Get(
			context.Background(),
			fmt.Sprintf("%v/%v", srvApi, "test_get_params"),
			nil,
			DefaultJsonHeader(),
		).BodyString()
		assert.Nil(t, err)
		assert.Equal(t, `{"message":"hello superwhys"}`, resp)
	})
}

func TestClientGet(t *testing.T) {
	t.Run("testContextFetch", func(t *testing.T) {
		resp, err := Cli.Get(
			context.Background(),
			fmt.Sprintf("%v/%v", srvApi, "test_get"),
			nil,
			DefaultJsonHeader(),
		).BodyString()
		assert.Nil(t, err)
		assert.Equal(t, `{"message":"do success"}`, resp)
	})
}

type message struct {
	Message string `json:"message"`
}

func TestClientGetWithCallBack(t *testing.T) {
	msg := &message{}

	t.Run("testContextFetch", func(t *testing.T) {
		resp := Cli.Get(
			context.Background(),
			fmt.Sprintf("%v/%v", srvApi, "test_get"),
			nil,
			DefaultJsonHeader(),
			func(c *Context) {
				defer c.Response.Body.Close()
				b, err := ioutil.ReadAll(c.Response.Body)
				if err != nil {
					fmt.Printf("read body error:%v\n", err)
					return
				}
				fmt.Printf("body:%v\n", string(b))
				err = json.Unmarshal(b, &msg)
				if err != nil {
					fmt.Printf("unmarshal body error:%v\n", err)
					return
				}
				fmt.Printf("message:%v\n", msg.Message)
			},
		)
		assert.Nil(t, resp.Error())
		assert.Equal(t, &message{Message: "do success"}, msg)
	})
}

func TestClientGetWithCallBack2(t *testing.T) {
	msg := &message{}

	t.Run("testContextFetch", func(t *testing.T) {
		resp := Cli.Get(
			context.Background(),
			fmt.Sprintf("%v/%v", srvApi, "test_get"),
			nil,
			DefaultJsonHeader(),
			func(c *Context) {
				if c.Response.StatusCode == 200 {
					c.Abort()
				}
			},
			func(c *Context) {
				defer c.Response.Body.Close()
				b, err := ioutil.ReadAll(c.Response.Body)
				if err != nil {
					fmt.Printf("read body error:%v\n", err)
					return
				}
				fmt.Printf("body:%v\n", string(b))
				err = json.Unmarshal(b, &msg)
				if err != nil {
					fmt.Printf("unmarshal body error:%v\n", err)
					return
				}
				fmt.Printf("message:%v\n", msg.Message)
			},
		)
		assert.Nil(t, resp.Error())
		assert.Equal(t, &message{}, msg)
	})
}
