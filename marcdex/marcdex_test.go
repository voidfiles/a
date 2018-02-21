package marcdex_test

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/boutros/marc"
	"github.com/voidfiles/a/marcdex"
	"github.com/voidfiles/a/recordstore"
)

type TestMarcStream struct {
	completed  int
	iterations int
}

func NewMarcStream(iterations int) *TestMarcStream {
	return &TestMarcStream{
		completed:  0,
		iterations: iterations,
	}
}

func (ms *TestMarcStream) Next() bool {
	if ms.completed >= ms.iterations {
		return false
	}
	ms.completed++
	return true
}
func (ms *TestMarcStream) Value() []*marc.Record {
	record := buildTestRecord()
	return []*marc.Record{
		&record,
	}
}

type TestNodeInterface struct{}

func (tn *TestNodeInterface) Save(interface{}) error {
	return nil
}

type TestDataWriter struct {
	err error
}

func (dw *TestDataWriter) SaveChunk(records []recordstore.ResoRecord) error { return dw.err }

func NewTestDataWriter(err error) *TestDataWriter {
	return &TestDataWriter{
		err: err,
	}
}
func TestMustNewMarcIndexer(t *testing.T) {
	ms := NewMarcStream(2)
	db := NewTestDataWriter(nil)
	marcdex.MustNewMarcIndexer(ms, db)
}

func TestBatchWrite(t *testing.T) {
	ms := NewMarcStream(1)
	db := NewTestDataWriter(nil)
	md := marcdex.MustNewMarcIndexer(ms, db)

	md.BatchWrite()
}

func TestBatchErr(t *testing.T) {
	ms := NewMarcStream(1)
	db := NewTestDataWriter(fmt.Errorf("Failed to write"))
	md := marcdex.MustNewMarcIndexer(ms, db)

	md.BatchWrite()
}

func buildTestRecord() marc.Record {
	subFieldA := marc.SubField{
		Code:  "a",
		Value: "xxx",
	}

	subFields := marc.SubFields{subFieldA}

	dField := marc.DField{
		Tag:       "155",
		Ind1:      "",
		Ind2:      "",
		SubFields: subFields,
	}

	cField := marc.CField{
		Tag:   "001",
		Value: "bbb",
	}

	return marc.Record{
		XMLName:    xml.Name{},
		Leader:     "",
		CtrlFields: marc.CFields{cField},
		DataFields: marc.DFields{dField},
	}
}
