package marcdex_test

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/voidfiles/a/marcdex"
)

func TestMustNewMarcIndexer(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fail()
	}
	tmpPath := tmpDir + "/db"

	marcdex.MustNewMarcIndexer(tmpPath)
	assert.Panics(t, func() {
		marcdex.MustNewMarcIndexer("/asdf/sdfsd/sdfs/sdfs")
	})

}
