package httpclient

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

var headers http.Header
var someErr error

func init() {
	headers := make(http.Header)
	headers.Add("key", "value")

	someErr = errors.New("some error")
}

func newClient(t *testing.T, opts ...Option) (Client, *MockDoer, func()) {
	ctrl := gomock.NewController(t)
	doer := NewMockDoer(ctrl)
	opts = append(opts, WithDoer(doer))
	cli, err := New(opts...)
	assert.Nil(t, err)
	return cli, doer, func() {
		ctrl.Finish()
	}
}

func TestHttpClient_DoSuccess(t *testing.T) {
	client, doer, done := newClient(t)
	defer done()

	// returns success (http status code 200)
	req, err := http.NewRequest(http.MethodPost, "https://google.com", nil)
	assert.Nil(t, err)
	wantResp := &http.Response{
		StatusCode: 200,
	}
	doer.EXPECT().Do(req).Times(1).Return(wantResp, nil)
	haveResp, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, wantResp, haveResp)

	// returns success (http status code 200 with request body)
	payload := []byte(`{"org": "mediabuyerbot"}`)
	req, err = http.NewRequest(http.MethodPost, "https://google.com", bytes.NewBuffer(payload))
	assert.Nil(t, err)
	wantResp = &http.Response{
		StatusCode: 200,
	}
	doer.EXPECT().Do(req).Times(1).Return(wantResp, nil).Do(func(r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		assert.Equal(t, payload, body)
	})
	haveResp, err = client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, wantResp, haveResp)
}

func TestHttpClient_DoError(t *testing.T) {
	client, doer, done := newClient(t)
	defer done()
	req, err := http.NewRequest(http.MethodGet, "https://google.com", nil)
	assert.Nil(t, err)
	doer.EXPECT().Do(req).Times(1).Return(nil, errors.New("refused"))
	haveResp, err := client.Do(req)
	assert.Error(t, err)
	assert.Nil(t, haveResp)
}

func TestHttpClient_DoErrorWithRetry(t *testing.T) {
	var (
		retryCount = 2
		haveRetry  = 0
		respErr    = errors.New("refused")
		payload    = []byte(`{"test":"test"}`)
	)
	client, doer, done := newClient(t,
		WithRetryCount(retryCount),
		WithResponseHook(func(request *http.Request, response *http.Response) {
			assert.Nil(t, response)
			assert.Nil(t, request)
		}),
		WithErrorHook(func(req *http.Request, err error, retry int) {
			assert.Equal(t, err, respErr)
			b, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err)
			assert.Equal(t, b, payload)
			haveRetry++
		}),
	)
	defer done()
	req, err := http.NewRequest(http.MethodPost, "https://google.com", bytes.NewBuffer(payload))
	assert.Nil(t, err)
	doer.EXPECT().Do(req).Times(retryCount+1).Return(nil, respErr)
	haveResp, err := client.Do(req)
	assert.Error(t, err)
	assert.Nil(t, haveResp)
	assert.Equal(t, haveRetry-1, retryCount)
}

func TestHttpClient_DoErrorWithRetryAndCheckRetry(t *testing.T) {
	var (
		retryCount = 3
		haveRetry  = 0
		respErr    = errors.New("refused")
		payload    = []byte(`{"test":"test"}`)
	)
	client, doer, done := newClient(t,
		WithRetryCount(retryCount),
		WithCheckRetry(func(req *http.Request, resp *http.Response, err error) (b bool, err2 error) {
			return false, errors.New("check err")
		}),
		WithErrorHook(func(req *http.Request, err error, retry int) {
			assert.Equal(t, err, respErr)
			b, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err)
			assert.Equal(t, b, payload)
			haveRetry++
		}),
	)
	defer done()
	req, err := http.NewRequest(http.MethodPost, "https://google.com", bytes.NewBuffer(payload))
	assert.Nil(t, err)
	doer.EXPECT().Do(req).Times(1).Return(nil, respErr)
	haveResp, err := client.Do(req)
	assert.Error(t, err)
	assert.Nil(t, haveResp)
	assert.Equal(t, haveRetry, 1)
}

func TestHttpClient_DoWithRetryHTTP500(t *testing.T) {
	var (
		payload     = []byte(`{"test":"test"}`)
		haveRetries int
		wantRetries = 6
	)
	client, doer, done := newClient(t,
		WithRetryCount(wantRetries-1),
		WithRequestHook(func(request *http.Request, i int) {
			haveRetries++
		}),
		WithResponseHook(func(request *http.Request, response *http.Response) {
			assert.Equal(t, response.StatusCode, http.StatusInternalServerError)
		}),
	)
	defer done()
	req, err := http.NewRequest(http.MethodPost, "https://google.com", bytes.NewBuffer(payload))
	assert.Nil(t, err)
	doer.EXPECT().Do(req).Times(wantRetries).Return(&http.Response{
		StatusCode: 500,
	}, nil).Do(func(r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		assert.Equal(t, payload, body)
	})
	haveResp, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, wantRetries, haveRetries)
	assert.Equal(t, haveResp.StatusCode, http.StatusInternalServerError)
}

func TestHttpClient_DoWithRetryHTTP200(t *testing.T) {
	var (
		payload     = []byte(`{"test":"test"}`)
		haveRetries int
		wantRetries = 1
	)
	client, doer, done := newClient(t,
		WithRetryCount(10),
		WithRequestHook(func(request *http.Request, i int) {
			haveRetries++
		}),
		WithResponseHook(func(request *http.Request, response *http.Response) {
			assert.Equal(t, response.StatusCode, http.StatusOK)
		}),
	)
	defer done()
	req, err := http.NewRequest(http.MethodPost, "https://google.com", bytes.NewBuffer(payload))
	assert.Nil(t, err)
	doer.EXPECT().Do(req).Times(wantRetries).Return(&http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(payload)),
	}, nil).Do(func(r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		assert.Equal(t, payload, body)
	})
	haveResp, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, wantRetries, haveRetries)
	assert.Equal(t, haveResp.StatusCode, http.StatusOK)

	b, err := ioutil.ReadAll(haveResp.Body)
	assert.Nil(t, err)
	assert.Equal(t, b, payload)
}

func TestHttpClient_DoWithRetryAndCheckRetryPolicyHTTP200(t *testing.T) {
	payload := []byte(`{"test":"test"}`)
	var haveRetries int
	client, doer, done := newClient(t,
		WithRetryCount(10),
		WithRequestHook(func(request *http.Request, i int) {
			haveRetries++
		}),
		WithErrorHandler(func(resp *http.Response, err error, numTries int) (response *http.Response, err2 error) {
			assert.Nil(t, err)
			assert.Equal(t, numTries, 0)
			return resp, nil
		}),
		WithCheckRetry(func(req *http.Request, resp *http.Response, err error) (b bool, err2 error) {
			if err != nil {
				return true, nil
			}
			if resp.StatusCode >= http.StatusInternalServerError {
				return true, nil
			}
			return false, nil
		}),
	)
	defer done()
	req, err := http.NewRequest(http.MethodPost, "https://google.com", bytes.NewBuffer(payload))
	assert.Nil(t, err)
	doer.EXPECT().Do(req).Times(1).Return(&http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(payload)),
	}, nil).Do(func(r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		assert.Equal(t, payload, body)
	})
	haveResp, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, 1, haveRetries)
	assert.Equal(t, haveResp.StatusCode, http.StatusOK)

	b, err := ioutil.ReadAll(haveResp.Body)
	assert.Nil(t, err)
	assert.Equal(t, b, payload)
}

func TestHttpClient_Delete(t *testing.T) {
	client, doer, done := newClient(t, WithBaseURL("http://test.com"))
	defer done()

	ctx := context.TODO()

	// returns success
	doer.EXPECT().Do(gomock.Any()).Times(1).Return(&http.Response{
		StatusCode: 200,
	}, nil).Do(func(req *http.Request) {
		assert.Equal(t, req.URL.Path, "/path")
		assert.Equal(t, req.URL.Host, "test.com")
		assert.Equal(t, req.Method, http.MethodDelete)
		assert.Equal(t, req.Header.Get("key"), headers.Get("key"))
	})

	resp, err := client.Delete(ctx, "/path", headers)
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	// returns error
	doer.EXPECT().Do(gomock.Any()).Times(1).Return(nil, someErr)
	resp, err = client.Delete(ctx, "/path", headers)
	assert.EqualError(t, err, someErr.Error())
	assert.Nil(t, resp)
}

func TestHttpClient_Put(t *testing.T) {
	client, doer, done := newClient(t, WithBaseURL("http://test.com"))
	defer done()

	payload := []byte(`{"test":"test"}`)
	ctx := context.TODO()

	// returns success
	doer.EXPECT().Do(gomock.Any()).Times(1).Return(&http.Response{
		StatusCode: 200,
	}, nil).Do(func(req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		assert.Nil(t, err)
		assert.Equal(t, req.URL.Path, "/path")
		assert.Equal(t, req.URL.Host, "test.com")
		assert.Equal(t, req.Method, http.MethodPut)
		assert.Equal(t, req.Header.Get("key"), headers.Get("key"))
		assert.Equal(t, body, payload)
	})

	resp, err := client.Put(ctx, "/path", ioutil.NopCloser(bytes.NewReader(payload)), headers)
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	// returns error
	doer.EXPECT().Do(gomock.Any()).Times(1).Return(nil, someErr)
	resp, err = client.Put(ctx, "/path", nil, headers)
	assert.EqualError(t, err, someErr.Error())
	assert.Nil(t, resp)
}

func TestHttpClient_Post(t *testing.T) {
	client, doer, done := newClient(t, WithBaseURL("http://test.com"))
	defer done()

	payload := []byte(`{"test":"test"}`)
	ctx := context.TODO()

	// returns success
	doer.EXPECT().Do(gomock.Any()).Times(1).Return(&http.Response{
		StatusCode: 200,
	}, nil).Do(func(req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		assert.Nil(t, err)
		assert.Equal(t, req.URL.Path, "/path")
		assert.Equal(t, req.URL.Host, "test.com")
		assert.Equal(t, req.Method, http.MethodPost)
		assert.Equal(t, req.Header.Get("key"), headers.Get("key"))
		assert.Equal(t, body, payload)
	})

	resp, err := client.Post(ctx, "/path", ioutil.NopCloser(bytes.NewReader(payload)), headers)
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	// returns error
	doer.EXPECT().Do(gomock.Any()).Times(1).Return(nil, someErr)
	resp, err = client.Post(ctx, "/path", nil, headers)
	assert.EqualError(t, err, someErr.Error())
	assert.Nil(t, resp)
}

func TestHttpClient_Get(t *testing.T) {
	client, doer, done := newClient(t, WithBaseURL("http://test.com"))
	defer done()

	ctx := context.TODO()

	// returns success
	doer.EXPECT().Do(gomock.Any()).Times(1).Return(&http.Response{
		StatusCode: 200,
	}, nil).Do(func(req *http.Request) {
		assert.Equal(t, req.URL.Path, "/path")
		assert.Equal(t, req.URL.Host, "test.com")
		assert.Equal(t, req.Method, http.MethodGet)
		assert.Equal(t, req.Header.Get("key"), headers.Get("key"))
	})

	resp, err := client.Get(ctx, "/path", headers)
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	// returns error
	doer.EXPECT().Do(gomock.Any()).Times(1).Return(nil, someErr)
	resp, err = client.Get(ctx, "/path", headers)
	assert.EqualError(t, err, someErr.Error())
	assert.Nil(t, resp)
}
