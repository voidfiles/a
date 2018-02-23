package recordstore

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/coreos/bbolt"
	"github.com/vmihailenco/msgpack"
)

const (
	// ResoRecordBucketName is the name of DB Buckte
	ResoRecordBucketName = "ResoRecord"
	// IdentifierKeyPrefix prefixes keys in boltdb for identifier field
	IdentifierKeyPrefix = "identifier"
	// AltIdentifierKeyPrefix prefixes keys in boltdb for alt-identifier field
	AltIdentifierKeyPrefix = "alt-identifier"
	// OldIdentifierKeyPrefix prefixes keys in boltdb for old-identifier field
	OldIdentifierKeyPrefix = "old-identifier"
)

// RecordStore will store a record into an index
type RecordStore struct {
	db *bolt.DB
}

// ResoRecord is a record we can use to do authority resolution
type ResoRecord struct {
	_msgpack        struct{} `msgpack:",omitempty"`
	Identifier      string   `json:"identifier,omitempty"`
	Type            string   `json:"type,omitempty"`
	AltIdentifier   []string `json:"alt-identifier,omitempty"`
	OldIdentifier   []string `json:"old-identifier,omitempty"`
	Heading         []string `json:"heading,omitempty"`
	AltHeading      []string `json:"alt-heading,omitempty"`
	WestCoordinate  []string `json:"west-coordinate,omitempty"`
	EastCoordinate  []string `json:"east-coordinate,omitempty"`
	NorthCoordinate []string `json:"north-coordinate,omitempty"`
	SouthCoordinate []string `json:"south-coordinate,omitempty"`
	MARCGeoCode     []string `json:"marc-geo-code,omitempty"`
	Classification  []string `json:"classification,omitempty"`
	GeneralNote     []string `json:"general-note,omitempty"`
}

// StorageOperation is an operation that will get stored in a key value database
type StorageOperation struct {
	Key    []byte
	Value  []byte
	Bucket string
}

// NewStorageOperation crates a new KeyValue
func NewStorageOperation(bucket string, keyPrefix string, key string, value []byte) StorageOperation {
	return StorageOperation{
		Key:    []byte(fmt.Sprintf("%s:%s", keyPrefix, key)),
		Value:  value,
		Bucket: bucket,
	}
}

// MustNewRecordStore will create a new RecordStore
func MustNewRecordStore(db *bolt.DB) *RecordStore {
	err := db.Batch(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(ResoRecordBucketName))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return &RecordStore{
		db: db,
	}
}

// ConvertRecordToStorageOperations returns a list of operations to be stored
func ConvertRecordToStorageOperations(record ResoRecord) ([]StorageOperation, error) {
	keyValues := make([]StorageOperation, 0)
	mainValue, err := msgpack.Marshal(record)

	if err != nil {
		return keyValues, nil
	}

	keyValues = append(keyValues, NewStorageOperation(ResoRecordBucketName, IdentifierKeyPrefix, record.Identifier, mainValue))

	for _, id := range record.AltIdentifier {
		keyValues = append(
			keyValues,
			NewStorageOperation(
				ResoRecordBucketName,
				fmt.Sprintf("%s:%s", AltIdentifierKeyPrefix, record.Identifier),
				id,
				[]byte(record.Identifier),
			),
		)
	}
	//
	for _, id := range record.OldIdentifier {
		keyValues = append(
			keyValues,
			NewStorageOperation(
				ResoRecordBucketName,
				fmt.Sprintf("%s:%s", OldIdentifierKeyPrefix, record.Identifier),
				id,
				[]byte(record.Identifier),
			),
		)
	}

	return keyValues, nil
}

//HandleOperation will persist an operation into the database
func HandleOperation(tx *bolt.Tx, operation StorageOperation) error {
	bucket := tx.Bucket([]byte(operation.Bucket))
	val := bucket.Get(operation.Key)
	if val != nil {
		log.Printf("While fetching key %s found exisisting %s", string(operation.Key), string(operation.Value))
	}
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
			operations, err := ConvertRecordToStorageOperations(record)
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

// FindByIdentifier will lookup a ResoRecord by its main identifier
func (r *RecordStore) FindByIdentifier(id string) (*ResoRecord, error) {
	var record ResoRecord
	err := r.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(ResoRecordBucketName))
		value := bucket.Get([]byte(fmt.Sprintf("%s:%s", IdentifierKeyPrefix, id)))

		msgpack.Unmarshal(value, &record)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (r *RecordStore) Stats() (string, error) {
	var statStr string
	err := r.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(ResoRecordBucketName))
		stats := bucket.Stats()
		statBytes, err := json.Marshal(stats)
		if err != nil {
			return err
		}
		statStr = string(statBytes)
		return nil
	})

	if err != nil {
		return "", err
	}

	return statStr, nil
}
