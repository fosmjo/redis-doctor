package doctor

import "time"

type Entry interface {
	Accept(Visitor) error
}

type Visitor interface {
	VisitBigKey(key *BigKey) error
	VisitSlowLog(log *SlowLog) error
}

type BigKey struct {
	Key              string `json:"key"`
	Type             string `json:"type"`
	Encoding         string `json:"encoding"`
	SerializedLength int64  `json:"serializedlength"`
	// Number of elements of the key, acquired by following commands:
	//   STRLEN <string>
	//   HLEN <hash>
	//   LLEN <list>
	//   SCARD <set>
	//   ZCARD <sorted set>
	Cardinality int64 `json:"cardinality"`
}

func (bk *BigKey) Accept(v Visitor) error {
	return v.VisitBigKey(bk)
}

// Copied from https://pkg.go.dev/github.com/redis/go-redis/v9@v9.0.2#SlowLog
type SlowLog struct {
	ID         int64         `json:"id,omitempty"`
	Time       time.Time     `json:"time,omitempty"`
	Duration   time.Duration `json:"duration,omitempty"`
	Args       []string      `json:"args,omitempty"`
	ClientAddr string        `json:"clientAddr,omitempty"`
	ClientName string        `json:"clientName,omitempty"`
}

func (sl *SlowLog) Accept(v Visitor) error {
	return v.VisitSlowLog(sl)
}
