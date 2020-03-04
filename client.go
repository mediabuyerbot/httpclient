package httpclient

import (
	"context"
	"io"
	"net/http"
	"time"
)

// Doer interface has the method required to use a type as custom http client.
// The net/*http.Client type satisfies this interface.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// Client is a generic HTTP client interface.
type Client interface {
	Get(ctx context.Context, url string, headers http.Header) (*http.Response, error)
	Post(ctx context.Context, url string, body io.Reader, headers http.Header) (*http.Response, error)
	Put(ctx context.Context, url string, body io.Reader, headers http.Header) (*http.Response, error)
	Delete(ctx context.Context, url string, headers http.Header) (*http.Response, error)
	Do(req *http.Request) (*http.Response, error)
}

// RequestHook allows a function to run before each retry. The HTTP
// request which will be made, and the retry number (0 for the initial
// request) are available to users.
type RequestHook func(*http.Request, int)

// ResponseHook is like RequestHook, but allows running a function
// on each HTTP response. This function will be invoked at the end of
// every HTTP request executed, regardless of whether a subsequent retry
// needs to be performed or not. If the response body is read or closed
// from this method, this will affect the response returned from Do().
type ResponseHook func(*http.Request, *http.Response)

// CheckRetry specifies a policy for handling retries. It is called
// following each request with the response and error values returned by
// the http.Client. If CheckRetry returns false, the Client stops retrying
// and returns the response to the caller. If CheckRetry returns an error,
// that error value is returned in lieu of the error from the request. The
// Client will close any response body when retrying, but if the retry is
// aborted it is up to the CheckRetry callback to properly close any
// response body before returning.
type CheckRetry func(req *http.Request, resp *http.Response, err error) (bool, error)

// BackOff specifies a policy for how long to wait between retries.
// It is called after a failing request to determine the amount of time
// that should pass before trying again.
type BackOff func(attemptNum int, resp *http.Response) time.Duration

// ErrorHandler is called if retries are expired, containing the last status
// from the http library. If not specified, default behavior for the library is
// to close the body and return an error indicating how many tries were
// attempted. If overriding this, be sure to close the body if needed.
type ErrorHandler func(resp *http.Response, err error, numTries int) (*http.Response, error)

// ErrorHook is called when the request returned a connection error.
type ErrorHook func(req *http.Request, err error, retry int)
