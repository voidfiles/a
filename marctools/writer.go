package marctools

import (
	"log"

	"github.com/voidfiles/a/pipeline"
	"github.com/voidfiles/a/recordstore"
)

type ResoRecordWriter struct {
	rs      *recordstore.RecordStore
	in      chan []recordstore.ResoRecord
	name    string
	written int64
}

func MustResoRecordWriter(rs *recordstore.RecordStore, in chan []recordstore.ResoRecord) *ResoRecordWriter {
	return &ResoRecordWriter{
		rs:      rs,
		in:      in,
		name:    "record-writer",
		written: 0,
	}
}

func (w *ResoRecordWriter) Write(killChan chan error) {
	for chunk := range w.in {
		err := w.rs.SaveChunk(chunk)
		if err != nil {
			log.Printf("Failed to write chunk to record store")
			killChan <- err
		}
		w.written = w.written + int64(len(chunk))
	}
}

func (w *ResoRecordWriter) Name() string {
	return w.name
}

func (w *ResoRecordWriter) Stats() pipeline.WriterStats {
	return pipeline.WriterStats{
		Written: w.written,
	}
}
