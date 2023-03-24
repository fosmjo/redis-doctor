package proto

import (
	"strings"

	"github.com/mitchellh/mapstructure"
)

type DebugObjectResult struct {
	RefCount         int64  `mapstructure:"refcount"`
	Encoding         string `mapstructure:"encoding"`
	SerializedLength int64  `mapstructure:"serializedlength"`
	LRU              int64  `mapstructure:"lru"`
	LRUSecondsIdle   int64  `mapstructure:"lru_seconds_idle"`
}

// debug object result example:
// Value at:0x60000377c070 refcount:1 encoding:listpack serializedlength:26 lru:1904357 lru_seconds_idle:6
func ParseDebugObjectResult(raw string) (result *DebugObjectResult, err error) {
	parts := strings.Split(raw, " ")[2:]
	m := make(map[string]string, len(parts))

	for i := 0; i < len(parts); i++ {
		kv := strings.Split(parts[i], ":")
		m[kv[0]] = kv[1]
	}

	result = &DebugObjectResult{}
	err = mapstructure.WeakDecode(m, result)
	if err != nil {
		return
	}

	return
}
