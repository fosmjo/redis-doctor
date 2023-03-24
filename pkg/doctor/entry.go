package doctor

import (
	"time"
)

type Entry interface {
	Accept(Visitor) error
}

type Visitor interface {
	VisitBigKey(key *BigKey) error
	VisitSlowLog(log *SlowLog) error
}

type BigKey struct {
	Key              string `json:"key,omitempty" xml:"key,omitempty"`
	Type             string `json:"type,omitempty" xml:"type,omitempty"`
	Encoding         string `json:"encoding,omitempty" xml:"encoding,omitempty"`
	SerializedLength int64  `json:"serializedLength,omitempty" xml:"serializedLength,omitempty"`
	// Number of elements of the key, acquired by following commands:
	//   STRLEN <string>
	//   HLEN <hash>
	//   LLEN <list>
	//   SCARD <set>
	//   ZCARD <sorted set>
	Cardinality int64 `json:"cardinality,omitempty" xml:"cardinality,omitempty"`
}

func (bk *BigKey) Accept(v Visitor) error {
	return v.VisitBigKey(bk)
}

// Copied from https://pkg.go.dev/github.com/redis/go-redis/v9@v9.0.2#SlowLog
type SlowLog struct {
	ID         int64         `json:"id,omitempty" xml:"id,omitempty"`
	Time       time.Time     `json:"time,omitempty" xml:"time,omitempty"`
	Duration   time.Duration `json:"duration,omitempty" xml:"duration,omitempty"`
	Args       []string      `json:"args,omitempty" xml:"args,omitempty"`
	ClientAddr string        `json:"clientAddr,omitempty" xml:"clientAddr,omitempty"`
	ClientName string        `json:"clientName,omitempty" xml:"clientName,omitempty"`
}

func (sl *SlowLog) Accept(v Visitor) error {
	return v.VisitSlowLog(sl)
}
