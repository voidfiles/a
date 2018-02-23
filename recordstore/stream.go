package recordstore

import (
	"github.com/coreos/bbolt"
)

// MarcStream is a marc record iterator
type RecordStream struct {
	store          *RecordStore
	bucket         *bolt.Bucket
	chunkSize      int
	currentResults []ResoRecord
	done           bool
	currentKey     []byte
}

// NewMarcStream creates and returns a MarcStream reader
func NewRecordStream(store *RecordStore, chunkSize int) (*RecordStream, error) {

	ms := &RecordStream{
		store:     store,
		chunkSize: chunkSize,
		done:      false,
	}

	return ms, nil
}

func (rs *RecordStream) Next() bool {
	if rs.done {
		return false
	}
	startingPrefix := []byte("")
	recordPage := rs.store.Scan(startingPrefix, rs.currentKey, 1000)
	if recordPage.More == false {
		rs.done = true
	}
	rs.currentResults = recordPage.Records
	rs.currentKey = recordPage.LastKey
	return true
}

func (ms *RecordStream) Value() []ResoRecord {
	return ms.currentResults
}
