package testtools

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/coreos/bbolt"
)

// tempfile returns a temporary file path.
func tempfile() string {
	f, err := ioutil.TempFile("", "bolt-")
	if err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
	if err := os.Remove(f.Name()); err != nil {
		panic(err)
	}
	return f.Name()
}

// NewTempBoltDB will create a temporary boltdb
func NewTempBoltDB() *bolt.DB {

	db, err := bolt.Open(tempfile(), 600, &bolt.Options{
		Timeout:        1 * time.Second,
		NoSync:         true,
		NoFreelistSync: true,
	})
	if err != nil {
		panic(err)
	}

	return db
}
