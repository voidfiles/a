package marcdex_test

import (
	"io/ioutil"
	"testing"

	"github.com/boutros/marc"
	"github.com/stretchr/testify/assert"
	"github.com/voidfiles/a/marcdex"
)

type TestMarcStream struct{}

func (ms *TestMarcStream) Next() bool            { return true }
func (ms *TestMarcStream) Value() []*marc.Record { return make([]*marc.Record, 0) }

func TestMustNewMarcIndexer(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fail()
	}
	tmpPath := tmpDir + "/db"
	ms := &TestMarcStream{}
	marcdex.MustNewMarcIndexer(tmpPath, ms)
	assert.Panics(t, func() {
		marcdex.MustNewMarcIndexer("/asdf/sdfsd/sdfs/sdfs", ms)
	})

}
