package marcdex

import (
	"log"
	"strings"
	"sync"

	"github.com/boutros/marc"
	"github.com/voidfiles/a/data_manager"
)

type ResoRecord struct {
	ID              string
	Type            string `storm:"index"`
	AltIdentifier   []string
	OldIdentifier   []string
	Heading         []string
	AltHeading      []string
	WestCoordinate  []string
	EastCoordinate  []string
	NorthCoordinate []string
	SouthCoordinate []string
	MARCGeoCode     []string
	Classification  []string
	GeneralNote     []string
}

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

func getFields(rec *marc.Record, tags []string, subfield string, ignores []string) []string {
	var Values []string

	for _, tag := range tags {

		Field := rec.GetDFields(tag)

		if len(subfield) > 0 {
			if len(Field) > 0 {
				for i := 0; i < len(Field); i++ {
					s := string(Field[i].SubField(subfield))
					Values = append(Values, s)
				}
			}
		} else {
			if len(Field) > 0 && len(ignores) == 0 {
				for i := 0; i < len(Field); i++ {
					s := ""
					for _, subf := range Field[i].SubFields {
						s += subf.Value + ". "
					}
					Values = append(Values, s)
				}
			} else {
				for i := 0; i < len(Field); i++ {
					s := ""
					for _, subf := range Field[i].SubFields {
						if !stringInSlice(subf.Code, ignores) {
							s += subf.Value + ". "
						}
					}
					Values = append(Values, s)
				}
			}
		}

	}

	return Values
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ConvertMarctoResoRecord(rec *marc.Record) ResoRecord {
	id, _ := rec.GetCField("001")
	return ResoRecord{
		ID:              strings.Replace(id.Value, " ", "", -1),
		AltIdentifier:   getFields(rec, []string{"010"}, "a", []string{}),
		OldIdentifier:   getFields(rec, []string{"010"}, "z", []string{}),
		Heading:         getFields(rec, []string{"100", "110", "111", "130", "150", "151", "155", "180", "181", "182", "185"}, "", []string{}),
		AltHeading:      getFields(rec, []string{"400", "500", "410", "510", "411", "430", "530", "450", "550", "451", "551", "455", "555", "480", "580", "581", "781", "482", "485", "585"}, "", []string{"w", "5"}),
		WestCoordinate:  getFields(rec, []string{"034"}, "d", []string{}),
		EastCoordinate:  getFields(rec, []string{"034"}, "e", []string{}),
		NorthCoordinate: getFields(rec, []string{"034"}, "f", []string{}),
		SouthCoordinate: getFields(rec, []string{"034"}, "g", []string{}),
		MARCGeoCode:     getFields(rec, []string{"043"}, "a", []string{}),
		Classification:  getFields(rec, []string{"050", "053", "072", "073"}, "", []string{}),
		GeneralNote:     getFields(rec, []string{"680"}, "", []string{}),
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
	records    []ResoRecord
	chunkIndex int
}

// BatchWrite will write marc data to boltdb
func (mi *MarcIndexer) BatchWrite() error {
	chunkChan := make(chan indexChunk, 2)
	go func() {
		chunks := 0
		for mi.ms.Next() {
			chunks++
			log.Printf("%v chunk read", chunks)
			records := make([]ResoRecord, 0)
			for _, record := range mi.ms.Value() {
				records = append(records, ConvertMarctoResoRecord(record))
			}
			chunkChan <- indexChunk{
				records:    records,
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
				for _, record := range chunk.records {
					mi.db.Save(&record)
				}
				// mi.db.InTransaction(func(dbx data_manager.NodeInterface) error {
				// 	for _, record := range chunk.records {
				// 		dbx.Save(&record)
				// 	}
				//
				// 	return nil
				// })

				log.Printf("%v chunk finished", chunk.chunkIndex)
			}
		}()
	}
	wg.Wait()
	return nil
}
