package proto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDebugObjectResult(t *testing.T) {
	raw := "Value at:0x60000377c070 refcount:1 encoding:listpack serializedlength:26 lru:1904357 lru_seconds_idle:6"
	result, err := ParseDebugObjectResult(raw)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.RefCount)
	assert.Equal(t, "listpack", result.Encoding)
	assert.Equal(t, int64(26), result.SerializedLength)
	assert.Equal(t, int64(1904357), result.LRU)
	assert.Equal(t, int64(6), result.LRUSecondsIdle)
}
