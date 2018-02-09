package authority

import (
	"fmt"
	"log"

	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"
	"github.com/voidfiles/a/search"
)

type Indexer struct {
	qs    graph.QuadStore
	index search.AuthoritySearch
}

func MustNewIndexer(qs graph.QuadStore, index search.AuthoritySearch) *Indexer {
	return &Indexer{
		qs,
		index,
	}
}

type QuadChunkStatefulIterator struct {
	qs        graph.QuadStore
	iterator  graph.Iterator
	quadChunk []search.QuadForIndex
	seen      int64
	totalSize int64
}

func NewQuadChunkStatefulIterator(qs graph.QuadStore) *QuadChunkStatefulIterator {
	iterator := qs.QuadsAllIterator()
	totalQuads, _ := iterator.Size()
	return &QuadChunkStatefulIterator{
		qs:        qs,
		iterator:  iterator,
		quadChunk: make([]search.QuadForIndex, 0),
		seen:      0,
		totalSize: totalQuads,
	}
}

func (it *QuadChunkStatefulIterator) Value() []search.QuadForIndex {
	return it.quadChunk
}

func (it *QuadChunkStatefulIterator) addQuadToChunk(val graph.Value, quad quad.Quad) {
	quadValue := search.NewQuadForIndex(
		fmt.Sprintf("%v", val),
		quad.Subject.String(),
		quad.Predicate.String(),
		quad.Object.String(),
	)
	it.quadChunk = append(it.quadChunk, quadValue)
}

func (it *QuadChunkStatefulIterator) Next() bool {
	it.quadChunk = make([]search.QuadForIndex, 0)

	for it.iterator.Next(nil) {
		it.seen++
		val := it.iterator.Result()
		quad := it.qs.Quad(val)

		if it.seen%10000 == 0 {
			done := (float64(it.seen) / float64(it.totalSize)) * 100
			log.Printf("%d of %d seen (%%%f done)", it.seen, it.totalSize, done)
		}

		if quad.Predicate.String() == "<http://www.w3.org/2004/02/skos/core#altLabel>" {
			it.addQuadToChunk(val, quad)
		}

		if quad.Predicate.String() == "<http://www.w3.org/2004/02/skos/core#prefLabel>" {
			it.addQuadToChunk(val, quad)
		}
		if len(it.quadChunk) > 100 {
			return true
		}

	}
	if len(it.quadChunk) > 0 {
		return true
	}

	return false
}

func (i *Indexer) chunkQuadStore(chunkChan chan []search.QuadForIndex) {
	chunker := NewQuadChunkStatefulIterator(i.qs)

	for chunker.Next() {
		quadChunk := chunker.Value()
		chunkChan <- quadChunk
	}

	close(chunkChan)
}

func (i *Indexer) indexChunks(chunkChan chan []search.QuadForIndex) {
	for quadChunk := range chunkChan {
		i.index.BatchIndex(quadChunk)
	}
}

func (i *Indexer) ProcessFullText() error {
	chunkChan := make(chan []search.QuadForIndex)
	go i.chunkQuadStore(chunkChan)
	i.indexChunks(chunkChan)

	return nil
}
