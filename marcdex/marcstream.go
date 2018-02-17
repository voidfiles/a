package marcdex

import (
	"io"
	"log"

	"github.com/boutros/marc"
)

// MarcStream is a marc record iterator
type MarcStream struct {
	decoder        *marc.Decoder
	chunkSize      int
	currentResults []*marc.Record
}

// NewMarcStream creates and returns a MarcStream reader
func NewMarcStream(data io.Reader, chunkSize int) (*MarcStream, error) {

	dec := marc.NewDecoder(data, marc.MARCXML)

	ms := &MarcStream{
		decoder:   dec,
		chunkSize: chunkSize,
	}

	return ms, nil
}

func (ms *MarcStream) Next() bool {
	ms.currentResults = make([]*marc.Record, 0)
	for i := 1; i <= ms.chunkSize; i++ {
		rec, err := ms.decoder.Decode()
		if err == io.EOF {
			return false
		}

		if err != nil {
			log.Printf("Failed to parse marc record %v", err)
		}
		ms.currentResults = append(ms.currentResults, rec)
	}
	return true
}

func (ms *MarcStream) Value() []*marc.Record {
	return ms.currentResults
}
