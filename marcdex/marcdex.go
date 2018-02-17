package marcdex

import (
	"github.com/boltdb/bolt"
	"github.com/boutros/marc"
)

// SubjectHeadingMarc is bilibographic record that gets indexed by id
type SubjectHeadingMarc struct {
	ID       string
	Headings []string
}

// MarcIndexer manages the marc index
type MarcIndexer struct {
	db *bolt.DB
	ms IMarcStream
}

type IMarcStream interface {
	Next() bool
	Value() []*marc.Record
}

// MustNewMarcIndexer will return a new MarcIndexer or die
func MustNewMarcIndexer(path string, ms IMarcStream) *MarcIndexer {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		panic(err)
	}

	m := &MarcIndexer{
		db: db,
		ms: ms,
	}

	return m
}
