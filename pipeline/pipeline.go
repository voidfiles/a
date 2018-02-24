package pipeline

import (
	"fmt"
	"sync"
	"time"
)

// ReaderStats we expect to see from a Reader
type ReaderStats struct {
	Read int64
}

// TransformStats we expect from a transformer
type TransformStats struct {
	Processed int64
}

// WriterStats we expect from a writer
type WriterStats struct {
	Written int64
}

// RecordReader is the interface a pipeline reader must adhere to
type RecordReader interface {
	Read(chan error)
	Finish()
	Stats() ReaderStats
	Name() string
}

// RecordTransform is the interface a pipeline transform must adhere to
type RecordTransform interface {
	Run(chan error)
	Stats() TransformStats
	Name() string
}

// RecordWriter is the interface a pipeline writer must adhere to
type RecordWriter interface {
	Write(chan error)
	Stats() WriterStats
	Name() string
}

// Pipeline manages an ETL pipeline
type Pipeline struct {
	reader     RecordReader
	writer     RecordWriter
	transforms []RecordTransform
	killChan   chan error
	wg         sync.WaitGroup
	name       string
	startTime  time.Time
}

// MustNewPipeline will create and return an ETL pipeline
func MustNewPipeline(name string, reader RecordReader, writer RecordWriter, transforms ...RecordTransform) *Pipeline {
	return &Pipeline{
		reader:     reader,
		writer:     writer,
		transforms: transforms,
		killChan:   make(chan error),
		name:       name,
	}
}

// Run will start the pipeline
func (p *Pipeline) Run() {
	p.startTime = time.Now()
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.reader.Read(p.killChan)
		p.reader.Finish()
	}()

	for _, transform := range p.transforms {
		p.wg.Add(1)
		go func(t RecordTransform) {
			defer p.wg.Done()
			t.Run(p.killChan)
		}(transform)
	}

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.writer.Write(p.killChan)
	}()

}

// Wait will wait for the pipeline to finish
func (p *Pipeline) Wait() {
	p.wg.Wait()
}

func (p *Pipeline) Stats() string {
	t := time.Now()
	elapsed := t.Sub(p.startTime)
	o := fmt.Sprintf("%s: %s\r\n", p.name, elapsed)
	o += fmt.Sprintf("Read (%s) %d)\r\n", p.reader.Name(), p.reader.Stats().Read)
	for _, transform := range p.transforms {
		o += fmt.Sprintf("Tran (%s) %d)\r\n", transform.Name(), transform.Stats().Processed)
	}
	o += fmt.Sprintf("Writ (%s) %d)\r\n", p.writer.Name(), p.writer.Stats().Written)
	return o
}
