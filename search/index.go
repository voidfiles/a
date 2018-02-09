package search

import (
	"fmt"
	"log"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/analysis/lang/en"
	"github.com/blevesearch/bleve/mapping"
)

// AuthoritySearch manages a search index
type AuthoritySearch interface {
	BatchIndex([]QuadForIndex) error
	Query(string) (*bleve.SearchResult, error)
}

type QuadForIndex struct {
	id        string
	subject   string
	predicate string
	object    string
}

func NewQuadForIndex(id, subject, predicate, object string) QuadForIndex {
	return QuadForIndex{
		id,
		subject,
		predicate,
		object,
	}
}

func (q QuadForIndex) Id() string {
	return q.id
}

func (q QuadForIndex) Subject() string {
	return q.subject
}

func (q QuadForIndex) Predicate() string {
	return q.predicate
}

func (q QuadForIndex) Object() string {
	return q.object
}

type SimpleIndex struct {
	Subject   string
	Predicate string
	Object    string
}

// Index is a struct that contains indexing info
type Index struct {
	index bleve.Index
}

// BatchIndex will index a bunch of QuadValues into the search index
func (i *Index) BatchIndex(quads []QuadForIndex) error {
	batch := i.index.NewBatch()
	for _, quad := range quads {
		quadForIndex := SimpleIndex{
			quad.Subject(),
			quad.Predicate(),
			quad.Object()}

		id := fmt.Sprintf("%v", quad.Id())

		batch.Index(id, quadForIndex)
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
	search.Fields = []string{"Object", "Subject", "Predicate"}
	return i.index.Search(search)
}

func buildQuadMapping() (mapping.IndexMapping, error) {
	// a generic reusable mapping for english text
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	// a generic reusable mapping for keyword text
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name

	quadMapping := bleve.NewDocumentMapping()

	// We want full text matching here
	quadMapping.AddFieldMappingsAt("Object", englishTextFieldMapping)

	// We want simple subject == "x" here
	quadMapping.AddFieldMappingsAt("Subject", keywordFieldMapping)
	quadMapping.AddFieldMappingsAt("Predicate", keywordFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("quad", quadMapping)

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
		indexMapping, berr := buildQuadMapping()
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
