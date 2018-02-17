package marcdex_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/voidfiles/a/marcdex"
)

func testMarcXml() (io.Reader, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("%s%s", wd, "/testdata/marc.xml")
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return fd, nil
}

func TestNewMarcStream(t *testing.T) {
	tmpPath, err := ioutil.TempFile("", "")
	assert.NoError(t, err)
	ms, err := marcdex.NewMarcStream(tmpPath, 100)
	assert.NoError(t, err)
	assert.IsType(t, &marcdex.MarcStream{}, ms)

}

func TestIterator(t *testing.T) {
	testData, err := testMarcXml()
	assert.NoError(t, err)
	ms, err := marcdex.NewMarcStream(testData, 2)
	assert.NoError(t, err)
	more := ms.Next()
	assert.True(t, more, "")
	records := ms.Value()
	assert.Len(t, records, 2)
	counter := 0
	for ms.Next() {
		counter++
	}

	assert.Equal(t, 5, counter)
}
