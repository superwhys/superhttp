package superhttp

type Params map[string]string

func (p Params) Get(key string) string {
	if val, exists := p[key]; !exists {
		return ""
	} else {
		return val
	}
}

func (p Params) Add(key, value string) Params {
	p[key] = value
	return p
}
