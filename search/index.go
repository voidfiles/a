package search

import (
	"fmt"
	"log"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/mapping"
	"github.com/voidfiles/a/recordstore"
)

type indexResoRecord struct {
	Heading    string
	AltHeading string
}

func unrollStr(strs []string) string {
	var start string
	for _, part := range strs {
		start += " " + part
	}

	return start
}

func ConvertResoRecordToIndexRecord(record recordstore.ResoRecord) indexResoRecord {
	return indexResoRecord{
		Heading:    unrollStr(record.Heading),
		AltHeading: unrollStr(record.AltHeading),
	}
}

// Index is a struct that contains indexing info
type Index struct {
	index bleve.Index
}

// BatchIndex will index a bunch of recordstore.ResoRecord into the search index
func (i *Index) BatchIndex(records []recordstore.ResoRecord) error {
	batch := i.index.NewBatch()
	for _, record := range records {
		indexRecord := ConvertResoRecordToIndexRecord(record)
		batch.Index(record.Identifier, indexRecord)
	}
	err := i.index.Batch(batch)
	if err != nil {
		return fmt.Errorf("BatchIndex failed: %s", err)
	}

	return nil
}

func (i *Index) Query(label string) (*bleve.SearchResult, error) {
	query := bleve.NewFuzzyQuery(label)
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{"Heading", "AltHeading"}
	return i.index.Search(search)
}

func buildRecordMapping() (mapping.IndexMapping, error) {
	// a generic reusable mapping for english text
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	resoMapping := bleve.NewDocumentMapping()

	// We want full text matching here
	resoMapping.AddFieldMappingsAt("Heading", englishTextFieldMapping)
	resoMapping.AddFieldMappingsAt("AltHeading", englishTextFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("record", resoMapping)

	indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping, nil
}

// MustNewIndex creates a new Index
func MustNewIndex(path string) *Index {

	index, err := bleve.Open(path)
	if err == bleve.ErrorIndexPathDoesNotExist {
		log.Printf("Creating new search index...")
		// create a mapping
		indexMapping, berr := buildRecordMapping()
		if berr != nil {
			log.Fatal(berr)
		}
		index, berr = bleve.New(path, indexMapping)
		if berr != nil {
			log.Fatal(berr)
		}
	} else if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Opening existing index...")
	}

	return &Index{
		index: index,
	}
}
