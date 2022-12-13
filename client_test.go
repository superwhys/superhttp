package superhttp

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var Cli *Client

func TestMain(m *testing.M) {
	Cli = Default()

	go func() {
		r := gin.Default()
		r.GET("/test_get", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "do success",
			})
		})
		panic(r.Run(":29940"))
	}()
	time.Sleep(2 * time.Second)

	os.Exit(m.Run())
}

func TestClientGet(t *testing.T) {
	t.Run("testContextFetch", func(t *testing.T) {
		resp, err := Cli.Get(
			context.Background(),
			fmt.Sprintf("%v/%v", srvApi, "test_get"),
			nil,
			DefaultJsonHeader(),
		).Body()
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
