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
	done           bool
}

// NewMarcStream creates and returns a MarcStream reader
func NewMarcStream(data io.Reader, chunkSize int, format marc.Format) (*MarcStream, error) {

	dec := marc.NewDecoder(data, format)

	ms := &MarcStream{
		decoder:   dec,
		chunkSize: chunkSize,
		done:      false,
	}

	return ms, nil
}

func (ms *MarcStream) Next() bool {
	if ms.done {
		return false
	}
	ms.currentResults = make([]*marc.Record, 0)
	for i := 1; i <= ms.chunkSize; i++ {
		rec, err := ms.decoder.Decode()
		if err == io.EOF {
			ms.done = true
			log.Printf("Reached end of file")
			return true
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
