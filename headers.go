package superhttp

import "net/http"

type Header struct {
	http.Header
}

func (h *Header) Add(key, value string) *Header {
	h.Set(key, value)
	return h
}

func DefaultJsonHeader() *Header {
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	return &Header{Header: header}
}

func DefaultFormUrlEncodedHeader() *Header {
	header := http.Header{}
	header.Add("Content-Type", "application/x-www-form-urlencoded")
	return &Header{Header: header}
}
