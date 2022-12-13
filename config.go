package superhttp

import "time"

type Config struct {
	requestTimeOut time.Duration
	proxy          string
}
