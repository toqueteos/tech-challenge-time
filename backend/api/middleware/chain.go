package middleware

import "net/http"

type Middleware func(prev http.Handler) http.Handler

type Chain struct {
	handlers []Middleware
}

func NewChain(mw ...Middleware) *Chain {
	return &Chain{}
}

func (c *Chain) Add(mw Middleware) {
	if mw != nil {
		c.handlers = append(c.handlers, mw)
	}
}

func (c *Chain) Apply(h http.Handler) http.Handler {
	for _, handler := range c.handlers {
		h = handler(h)
	}

	return h
}
