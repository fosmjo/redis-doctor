package visitors

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/fosmjo/redis-doctor/pkg/doctor"
)

type CSVVisitor struct {
	w *csv.Writer
}

func NewCSVVisitor(w io.Writer) *CSVVisitor {
	return &CSVVisitor{
		w: csv.NewWriter(w),
	}
}

func (v *CSVVisitor) VisitBigKey(key *doctor.BigKey) error {
	a := []string{
		key.Key,
		key.Type,
		key.Encoding,
		strconv.FormatInt(key.SerializedLength, 10),
		strconv.FormatInt(key.Cardinality, 10),
	}
	return v.w.Write(a)
}

func (cv *CSVVisitor) VisitSlowLog(log *doctor.SlowLog) error {
	a := []string{
		strconv.FormatInt(log.ID, 10),
		log.Time.Local().Format(time.DateTime),
		log.Duration.String(),
		strings.Join(log.Args, " "),
		log.ClientAddr,
		log.ClientName,
	}
	return cv.w.Write(a)
}

func (cv *CSVVisitor) Flush() {
	cv.w.Flush()
}
