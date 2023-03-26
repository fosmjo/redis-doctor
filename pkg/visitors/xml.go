package visitors

import (
	"encoding/xml"
	"io"

	"github.com/fosmjo/redis-doctor/pkg/doctor"
)

type XMLVisitor struct {
	w io.Writer
	e *xml.Encoder
}

func NewXMLVisitor(w io.Writer) *XMLVisitor {
	return &XMLVisitor{
		w: w,
		e: xml.NewEncoder(w),
	}
}

func (v *XMLVisitor) VisitBigKey(key *doctor.BigKey) error {
	t := struct {
		XMLName struct{} `xml:"bigkey"`
		*doctor.BigKey
	}{BigKey: key}

	return v.visit(&t)
}

func (v *XMLVisitor) VisitHotKey(key *doctor.HotKey) error {
	t := struct {
		XMLName struct{} `xml:"hotkey"`
		*doctor.HotKey
	}{HotKey: key}

	return v.visit(&t)
}

func (v *XMLVisitor) VisitSlowLog(log *doctor.SlowLog) error {
	t := struct {
		XMLName struct{} `xml:"slowlog"`
		*doctor.SlowLog
	}{SlowLog: log}

	return v.visit(&t)
}

func (v *XMLVisitor) visit(entry any) error {
	err := v.e.Encode(entry)
	if err != nil {
		return err
	}
	_, err = v.w.Write([]byte{'\n'})
	return err
}

func (v *XMLVisitor) Close() error {
	return v.e.Close()
}
