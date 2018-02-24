package lcsh

import (
	"github.com/boutros/marc"
	"github.com/voidfiles/a/marcdex"
	"github.com/voidfiles/a/pipeline"
	"github.com/voidfiles/a/recordstore"
)

// Transform is a pipeline transform step to
// convert marc.Record to recordstore.ResoRecord
type Transform struct {
	In        chan []*marc.Record
	Out       chan []recordstore.ResoRecord
	processed int64
	name      string
}

// MustNewTransform creates a lcsh Transform
func MustNewTransform(in chan []*marc.Record) *Transform {
	return &Transform{
		In:        in,
		Out:       make(chan []recordstore.ResoRecord),
		processed: 0,
		name:      "lcsh:marc2reso",
	}
}

// Run will start the pipeline process
func (t *Transform) Run(killChan chan error) {
	for item := range t.In {
		t.Out <- t.Transform(item)
	}
	close(t.Out)
}

// Transform will convert a chunk of marc.Records into a chunk of recordstore.ResoRecords
func (t *Transform) Transform(chunk []*marc.Record) []recordstore.ResoRecord {
	output := make([]recordstore.ResoRecord, 0)
	for _, record := range chunk {
		output = append(output, marcdex.ConvertMarctoResoRecord(record))
		t.processed++
	}

	return output
}

// Stats returns info about transform
func (t *Transform) Stats() pipeline.TransformStats {
	return pipeline.TransformStats{
		Processed: t.processed,
	}
}

func (t *Transform) Name() string {
	return t.name
}
