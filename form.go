package superhttp

import "net/url"

type Form struct {
	url.Values
}

func NewForm() *Form {
	v := url.Values{}
	return &Form{Values: v}
}

func (f *Form) Add(key, value string) *Form {
	f.Set(key, value)
	return f
}
