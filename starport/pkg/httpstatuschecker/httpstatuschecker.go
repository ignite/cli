// httpstatuschecker is a tool check health of http pages.
package httpstatuschecker

import (
	"net/http"
)

type checker struct {
	c      *http.Client
	addr   string
	method string
}

// Option used to customize checker.
type Option func(*checker)

// Client configures http client.
func Client(c *http.Client) Option {
	return func(cr *checker) {
		cr.c = c
	}
}

// Client configures http method.
func Method(name string) Option {
	return func(cr *checker) {
		cr.method = name
	}
}

// Check checks if given http addr is alive by applying options.
func Check(addr string, options ...Option) (isAvailable bool, err error) {
	cr := &checker{
		c:      http.DefaultClient,
		addr:   addr,
		method: http.MethodGet,
	}
	for _, o := range options {
		o(cr)
	}
	return cr.check()
}

func (c *checker) check() (bool, error) {
	req, err := http.NewRequest(c.method, c.addr, nil)
	if err != nil {
		return false, err
	}
	res, err := c.c.Do(req)
	if err != nil {
		return false, nil
	}
	defer res.Body.Close()
	isOKStatus := res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusMultipleChoices
	return isOKStatus, nil
}
