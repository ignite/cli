// Package httpstatuschecker is a tool check health of http pages.
package httpstatuschecker

import (
	"context"
	"net/http"
)

type checker struct {
	c      *http.Client
	addr   string
	method string
}

// Option used to customize checker.
type Option func(*checker)

// Method configures http method.
func Method(name string) Option {
	return func(cr *checker) {
		cr.method = name
	}
}

// Check checks if given http addr is alive by applying options.
func Check(ctx context.Context, addr string, options ...Option) (isAvailable bool, err error) {
	cr := &checker{
		c:      http.DefaultClient,
		addr:   addr,
		method: http.MethodGet,
	}
	for _, o := range options {
		o(cr)
	}
	return cr.check(ctx)
}

func (c *checker) check(ctx context.Context) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, c.method, c.addr, nil)
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
