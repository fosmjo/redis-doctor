package visitors

import (
	"encoding/json"
	"io"

	"github.com/fosmjo/redis-doctor/pkg/doctor"
)

type JSONVisitor struct {
	e *json.Encoder
}

func NewJSONVisitor(w io.Writer) *JSONVisitor {
	return &JSONVisitor{
		e: json.NewEncoder(w),
	}
}

func (v *JSONVisitor) VisitBigKey(key *doctor.BigKey) error {
	return v.e.Encode(key)
}

func (v *JSONVisitor) VisitSlowLog(log *doctor.SlowLog) error {
	return v.e.Encode(log)
}
