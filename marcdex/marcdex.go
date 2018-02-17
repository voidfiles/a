package marcdex

import (
	"github.com/boltdb/bolt"
)

// SubjectHeadingMarc is bilibographic record that gets indexed by id
type SubjectHeadingMarc struct {
	ID       string
	Headings []string
}

// MarcIndexer manages the marc index
type MarcIndexer struct {
	db *bolt.DB
}

// MustNewMarcIndexer will return a new MarcIndexer or die
func MustNewMarcIndexer(path string) *MarcIndexer {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		panic(err)
	}
	m := &MarcIndexer{
		db: db,
	}

	return m
}
