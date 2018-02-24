package marctools

import (
	"io"

	"github.com/boutros/marc"
	"github.com/voidfiles/a/marcdex"
	"github.com/voidfiles/a/pipeline"
)

type MarcReader struct {
	ms   *marcdex.MarcStream
	Out  chan []*marc.Record
	read int64
	name string
}

func MustNewMarcReader(data io.Reader, format marc.Format) *MarcReader {
	ms, err := marcdex.NewMarcStream(data, 1000, format)
	if err != nil {
		panic(err)
	}
	return &MarcReader{
		ms:   ms,
		Out:  make(chan []*marc.Record),
		read: 0,
		name: "marc-reader",
	}
}

func (r *MarcReader) Read(killChan chan error) {
	r.read = 0
	for r.ms.Next() {
		r.Out <- r.ms.Value()
		r.read = r.read + int64(len(r.ms.Value()))
	}
}

// Finish close out the channel
func (r *MarcReader) Finish() {
	close(r.Out)
}

// Stats Report back stats on the reader
func (r *MarcReader) Stats() pipeline.ReaderStats {
	return pipeline.ReaderStats{
		Read: r.read,
	}
}

func (r *MarcReader) Name() string {
	return r.name
}
