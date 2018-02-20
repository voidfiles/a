package marcdex

import (
	"log"
	"sync"

	"github.com/boutros/marc"
	"github.com/voidfiles/a/data_manager"
)

// SubjectHeadingMarc is bilibographic record that gets indexed by id
type SubjectHeadingMarc struct {
	ID       string
	Headings map[string]string
}

// ConvertMarcRecordToSubjectHeadingMarc converts a marc.Record to a SubjectHeadingMarc
func ConvertMarcRecordToSubjectHeadingMarc(m *marc.Record) SubjectHeadingMarc {
	cfield, _ := m.GetCField("001")
	id := cfield.Value
	headings := map[string]string{}
	for _, tag := range []string{"155", "455", "555"} {
		fields := m.GetDFields(tag)
		for _, field := range fields {
			heading := field.SubField("a")
			if heading != "" {
				headings[tag] = heading
			}
		}
	}
	return SubjectHeadingMarc{
		ID:       id,
		Headings: headings,
	}
}

// MarcIndexer manages the marc index
type MarcIndexer struct {
	db DataWriter
	ms IMarcStream
}

// IMarcStream is an interface for marc stream
type IMarcStream interface {
	Next() bool
	Value() []*marc.Record
}

// DataWriter is an interface that can save things to databases
// you can also save things in transactions
type DataWriter interface {
	Save(interface{}) error
	InTransaction(data_manager.TransactionFunction) error
}

// MustNewMarcIndexer will return a new MarcIndexer or die
func MustNewMarcIndexer(ms IMarcStream, db DataWriter) *MarcIndexer {
	m := &MarcIndexer{
		db: db,
		ms: ms,
	}

	return m
}

type indexChunk struct {
	headings   []SubjectHeadingMarc
	chunkIndex int
}

// BatchWrite will write marc data to boltdb
func (mi *MarcIndexer) BatchWrite() error {
	chunkChan := make(chan indexChunk, 4)
	go func() {
		chunks := 0
		for mi.ms.Next() {
			chunks++
			log.Printf("%v chunk read", chunks)
			headings := make([]SubjectHeadingMarc, 0)
			for _, record := range mi.ms.Value() {
				headings = append(headings, ConvertMarcRecordToSubjectHeadingMarc(record))
			}
			chunkChan <- indexChunk{
				headings:   headings,
				chunkIndex: chunks,
			}
		}
		close(chunkChan)
	}()
	var wg sync.WaitGroup
	for i := 0; i <= 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for chunk := range chunkChan {
				log.Printf("%v chunk working", chunk.chunkIndex)
				mi.db.InTransaction(func(dbx data_manager.NodeInterface) error {
					for _, heading := range chunk.headings {
						dbx.Save(heading)
					}

					return nil
				})
				log.Printf("%v chunk finished", chunk.chunkIndex)
			}
		}()
	}
	wg.Wait()
	return nil
}
