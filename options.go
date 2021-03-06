package httpclient

import (
	"time"
)

type Option func(client *HttpClient)

func WithDoer(client Doer) Option {
	return func(c *HttpClient) {
		if client == nil {
			return
		}
		c.client = client
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *HttpClient) {
		c.timeouts = timeout
	}
}

func WithRetryCount(retryCount int) Option {
	return func(c *HttpClient) {
		c.retryCount = retryCount
	}
}

func WithRequestHook(rh RequestHook) Option {
	return func(c *HttpClient) {
		c.requestHook = rh
	}
}

func WithResponseHook(rh ResponseHook) Option {
	return func(c *HttpClient) {
		c.responseHook = rh
	}
}

func WithCheckRetry(cr CheckRetry) Option {
	return func(c *HttpClient) {
		c.checkRetry = cr
	}
}

func WithBackOff(b BackOff) Option {
	return func(c *HttpClient) {
		if b == nil {
			return
		}
		c.backOff = b
	}
}

func WithErrorHook(eh ErrorHook) Option {
	return func(c *HttpClient) {
		c.errorHook = eh
	}
}

func WithErrorHandler(eh ErrorHandler) Option {
	return func(c *HttpClient) {
		c.errorHandler = eh
	}
}

func WithBaseURL(u string) Option {
	return func(c *HttpClient) {
		c.baseURL = u
	}
}
