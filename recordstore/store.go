package recordstore

import (
	"encoding/json"
	"fmt"

	"github.com/coreos/bbolt"
)

const dbResoRecord = "ResoRecord"

// RecordStore will store a record into an index
type RecordStore struct {
	db            *bolt.DB
	inTransaction bool
}

// ResoRecord is a record we can use to do authority resolution
type ResoRecord struct {
	Identifier      string   `json:"identifier"`
	Type            string   `json:"type"`
	AltIdentifier   []string `json:"alt-identifier"`
	OldIdentifier   []string `json:"old-identifier"`
	Heading         []string `json:"heading"`
	AltHeading      []string `json:"alt-heading"`
	WestCoordinate  []string `json:"west-coordinate"`
	EastCoordinate  []string `json:"east-coordinate"`
	NorthCoordinate []string `json:"north-coordinate"`
	SouthCoordinate []string `json:"south-coordinate"`
	MARCGeoCode     []string `json:"marc-geo-code"`
	Classification  []string `json:"classification"`
	GeneralNote     []string `json:"general-note"`
}

// KeyValue is an operation that will get stored in a key value database
type KeyValue struct {
	Key    []byte
	Value  []byte
	Bucket string
}

// NewKeyValue crates a new KeyValue
func NewKeyValue(bucket string, keyPrefix string, key string, value []byte) KeyValue {
	return KeyValue{
		Key:    []byte(fmt.Sprintf("%s:%s", keyPrefix, key)),
		Value:  value,
		Bucket: bucket,
	}
}

// MustNewRecordStore will create a new RecordStore
func MustNewRecordStore(db *bolt.DB) *RecordStore {
	err := db.Batch(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(dbResoRecord))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return &RecordStore{
		db:            db,
		inTransaction: false,
	}
}

// ConvertRecordToKeyValues returns a list of operations to be stored
func ConvertRecordToKeyValues(record ResoRecord) ([]KeyValue, error) {
	keyValues := make([]KeyValue, 0)
	mainValue, err := json.Marshal(record)
	if err != nil {
		return keyValues, nil
	}

	keyValues = append(keyValues, NewKeyValue(dbResoRecord, "primary_key", record.Identifier, mainValue))

	for _, id := range record.AltIdentifier {
		keyValues = append(
			keyValues,
			NewKeyValue(
				dbResoRecord,
				fmt.Sprintf("old-identifier:%s", record.Identifier),
				id,
				[]byte(record.Identifier),
			),
		)
	}
	//
	for _, id := range record.OldIdentifier {
		keyValues = append(
			keyValues,
			NewKeyValue(
				dbResoRecord,
				fmt.Sprintf("alt_key:%s", record.Identifier),
				id,
				[]byte(record.Identifier),
			),
		)
	}

	return keyValues, nil
}

//HandleOperation will persist an operation into the database
func HandleOperation(tx *bolt.Tx, operation KeyValue) error {
	bucket := tx.Bucket([]byte(operation.Bucket))
	err := bucket.Put(operation.Key, operation.Value)
	if err != nil {
		return err
	}

	return nil
}

// SaveChunk persists a chunk of ResoRecords to database
func (r *RecordStore) SaveChunk(records []ResoRecord) error {
	return r.db.Update(func(tx *bolt.Tx) error {

		for _, record := range records {
			operations, err := ConvertRecordToKeyValues(record)
			if err != nil {
				return fmt.Errorf("SaveChunk::ConvertRecordToKeyValues: %s", err)
			}
			for _, operation := range operations {
				err := HandleOperation(tx, operation)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}
