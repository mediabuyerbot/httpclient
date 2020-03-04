package httpclient

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gojek/valkyrie"

	"github.com/pkg/errors"
)

const (
	DefaultHTTPTimeout = 60 * time.Second
)

// HttpClient is the http client implementation
type HttpClient struct {
	baseURL      string
	client       Doer
	retryCount   int
	requestHook  RequestHook
	responseHook ResponseHook
	errorHook    ErrorHook
	checkRetry   CheckRetry
	backOff      BackOff
	errorHandler ErrorHandler
}

var defaultBackOffPolicy = func(attemptNum int, resp *http.Response) time.Duration {
	return 500 * time.Millisecond
}

// New returns a new instance of Client.
func New(opts ...Option) (Client, error) {
	client := HttpClient{
		backOff: defaultBackOffPolicy,
		client: &http.Client{
			Timeout: DefaultHTTPTimeout,
		},
	}
	for _, opt := range opts {
		opt(&client)
	}
	return &client, nil
}

// Get makes a HTTP GET request to provided URL.
func (c *HttpClient) Get(ctx context.Context, url string, headers http.Header) (*http.Response, error) {
	var response *http.Response
	if len(c.baseURL) > 0 {
		url = c.baseURL + url
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return response, errors.Wrap(err, "GET - request creation failed")
	}
	request.Header = headers
	return c.Do(request)
}

// Post makes a HTTP POST request to provided URL and requestBody.
func (c *HttpClient) Post(ctx context.Context, url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	if len(c.baseURL) > 0 {
		url = c.baseURL + url
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return response, errors.Wrap(err, "POST - request creation failed")
	}

	request.Header = headers

	return c.Do(request)
}

// Put makes a HTTP PUT request to provided URL and requestBody.
func (c *HttpClient) Put(ctx context.Context, url string, body io.Reader, headers http.Header) (*http.Response, error) {
	var response *http.Response
	if len(c.baseURL) > 0 {
		url = c.baseURL + url
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPut, url, body)
	if err != nil {
		return response, errors.Wrap(err, "PUT - request creation failed")
	}

	request.Header = headers

	return c.Do(request)
}

// Delete makes a HTTP DELETE request with provided URL.
func (c *HttpClient) Delete(ctx context.Context, url string, headers http.Header) (*http.Response, error) {
	var response *http.Response
	if len(c.baseURL) > 0 {
		url = c.baseURL + url
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return response, errors.Wrap(err, "DELETE - request creation failed")
	}

	request.Header = headers

	return c.Do(request)
}

// Do makes an HTTP request with the native `http.Do` interface.
func (c *HttpClient) Do(req *http.Request) (resp *http.Response, err error) {
	var bodyReader *bytes.Reader

	req.Close = true
	if req.Body != nil {
		reqData, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(reqData)
		req.Body = ioutil.NopCloser(bodyReader)
	}

	multiErr := &valkyrie.MultiError{}
	var numTries int
	for i := 0; i <= c.retryCount; i++ {
		isRetryOk := c.retryCount > 0 && i < c.retryCount
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}

		if c.requestHook != nil {
			c.requestHook(req, i)
		}

		var err error
		resp, err = c.client.Do(req)
		if bodyReader != nil {
			_, _ = bodyReader.Seek(0, 0)
		}
		if err != nil {
			if c.errorHook != nil {
				c.errorHook(req, err, i)
			}

			multiErr.Push(err.Error())

			if c.checkRetry != nil {
				checkOK, checkErr := c.checkRetry(req, resp, err)
				if !checkOK {
					if checkErr != nil {
						multiErr.Push(checkErr.Error())
					}
					break
				}
			}
			if isRetryOk {
				wait := c.backOff(i, resp)
				time.Sleep(wait)
			}
			numTries++
			continue
		}

		if c.responseHook != nil {
			c.responseHook(req, resp)
		}

		var nextLoop bool
		isDefaultRetryPolicy := resp.StatusCode >= http.StatusInternalServerError && isRetryOk

		if c.checkRetry != nil && isRetryOk {
			checkOK, checkErr := c.checkRetry(req, resp, nil)
			if !checkOK {
				if checkErr != nil {
					multiErr.Push(checkErr.Error())
				}
				break
			}
			nextLoop = true
		} else if isDefaultRetryPolicy {
			nextLoop = true
		}

		if nextLoop {
			wait := c.backOff(i, resp)
			time.Sleep(wait)
			numTries++
			continue
		}
		break
	}
	if c.errorHandler != nil {
		return c.errorHandler(resp, multiErr.HasError(), numTries)
	}
	return resp, multiErr.HasError()
}
