package superhttp

import (
	"io/ioutil"
	"net/http"
)

type Response struct {
	*http.Response
	err error
}

func (r *Response) Error() error {
	return r.err
}

func (r *Response) Body() (string, error) {
	bytes, err := r.BodyBytes()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (r *Response) BodyBytes() ([]byte, error) {
	if r.err != nil {
		return []byte{}, r.err
	}

	defer r.Response.Body.Close()
	bytes, err := ioutil.ReadAll(r.Response.Body)
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}
