package httpclient

import (
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	cli, err := New(
		WithTimeout(time.Second),
		WithBackOff(defaultBackOffPolicy),
	)
	httpcli, ok := cli.(*HttpClient)
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Equal(t, httpcli.timeouts, time.Second)
	assert.Equal(t,
		runtime.FuncForPC(reflect.ValueOf(httpcli.backOff).Pointer()).Name(),
		runtime.FuncForPC(reflect.ValueOf(defaultBackOffPolicy).Pointer()).Name(),
	)
}
