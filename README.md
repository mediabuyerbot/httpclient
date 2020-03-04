# httpclient  [![Coverage Status](https://coveralls.io/repos/github/mediabuyerbot/httpclient/badge.svg?branch=master)](https://coveralls.io/github/mediabuyerbot/httpclient?branch=master)

## Table of contents
- [Installation](#installation)
- [Commands](#commands)
  + [Build dependencies](#build-dependencies)
  + [Run test](#run-test)
  + [Run test with coverage profile](#run-test-with-coverage-profile) 
  + [Run sync coveralls](#run-sync-coveralls)
  + [Build mocks](#build-mocks) 
- [Usage](#usage)
  + [Making a GET request](#making-a-get-request)
  + [Making a GET request with headers](#making-a-get-request-with-headers) 
  + [Making a POST request](#making-a-post-request)
  + [Making a POST request with headers](#making-a-post-request-with-headers)
  + [Making a PUT request](#making-a-put-request)
  + [Making a PUT request with headers](#making-a-put-request-with-headers)
  + [Making a DELETE request](#making-a-delete-request)
  + [Making a DELETE request with headers](#making-a-delete-request-with-headers)
  + [Making a CUSTOM request](#making-a-custom-request)
- [Options](#options)
     
### Installation
```shell script
go get github.com/mediabuyerbot/httpclient
```

### Commands
#### Build dependencies
```shell script
make deps
```

#### Run test
```shell script
make test
```

#### Run test with coverage profile
```shell script
make covertest
```
#### Run sync coveralls
```shell script
COVERALLS_HTTPCLIENT_TOKEN=${COVERALLS_REPO_TOKEN}
make sync-coveralls
```
#### Build mocks
```shell script
make mocks
```

### Usage
#### Making a GET request 
```go
cli, err := New()
if err != nil {
    panic(err)
}
resp, err := cli.Get(context.TODO(), "https://google.com", nil)
if err != nil {
    panic(err)
}
...
```

#### Making a GET request with headers
```go
cli, err := New()
if err != nil {
    panic(err)
}
headers := make(http.Header)
headers.Add("X-Header", "value")
resp, err := cli.Get(context.TODO(), "https://google.com", headers)
if err != nil {
    panic(err)
}
...
```

#### Making a POST request 
```go
cli, err := New()
if err != nil {
    panic(err)
}
payload := ioutil.NopCloser(bytes.NewReader([]byte(`payload`)))
resp, err := cli.Post(context.TODO(), payload, "https://google.com", nil)
if err != nil {
    panic(err)
}
...
```

#### Making a POST request with headers
```go
cli, err := New()
if err != nil {
    panic(err)
}
headers := make(http.Header)
headers.Add("X-Header", "value")
payload := ioutil.NopCloser(bytes.NewReader([]byte(`payload`)))
resp, err := cli.POST(context.TODO(), payload, "https://google.com", headers)
if err != nil {
    panic(err)
}
...
```

#### Making a PUT request 
```go
cli, err := New()
if err != nil {
    panic(err)
}
payload := ioutil.NopCloser(bytes.NewReader([]byte(`payload`)))
resp, err := cli.PUT(context.TODO(), payload, "https://google.com", nil)
if err != nil {
    panic(err)
}
...
```

#### Making a PUT request with headers
```go
cli, err := New()
if err != nil {
    panic(err)
}
headers := make(http.Header)
headers.Add("X-Header", "value")
payload := ioutil.NopCloser(bytes.NewReader([]byte(`payload`)))
resp, err := cli.PUT(context.TODO(), payload, "https://google.com", headers)
if err != nil {
    panic(err)
}
...
```

#### Making a DELETE request 
```go
cli, err := New()
if err != nil {
    panic(err)
}
resp, err := cli.Delete(context.TODO(), "https://google.com", nil)
if err != nil {
    panic(err)
}
...
```

#### Making a DELETE request with headers
```go
cli, err := New()
if err != nil {
    panic(err)
}
headers := make(http.Header)
headers.Add("X-Header", "value")
resp, err := cli.Delete(context.TODO(), "https://google.com", headers)
if err != nil {
    panic(err)
}
...
```

#### Making a CUSTOM request
```go
cli, err := New()
if err != nil {
    panic(err)
}
req, err := http.NewRequest(http.MethodHead, "https://google.com", nil)
if err != nil {
    panic(err) 
}
resp, err := cli.Do(req)
if err != nil {
    panic(err)
}
...
``` 

### Options
```go
_, err := New(
   WithTimeout(time.Second),
   WithRetryCount(2),
   WithRequestHook(func(request *http.Request, i int) {})   
   WithResponseHook(func(request *http.Request, response *http.Response) {}),
   WithCheckRetry(func(req *http.Request, resp *http.Response, err error) (bool, error) {}),
   WithBackOff(func(attemptNum int, resp *http.Response) time.Duration {}),
   WithErrorHook(func(req *http.Request, err error, retry int) {}),
   WithErrorHandler(func(resp *http.Response, err error, numTries int) (*http.Response, error) {}),
   WithBaseURL("http://127.0.0.1"),  
)
```
