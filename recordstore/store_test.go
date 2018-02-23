package recordstore_test

import (
	"fmt"
	"testing"

	"github.com/coreos/bbolt"
	"github.com/stretchr/testify/assert"
	"github.com/voidfiles/a/recordstore"
	"github.com/voidfiles/a/testtools"
)

func TestNewStorageOperation(t *testing.T) {

	op := recordstore.NewStorageOperation("a", "b", "c", []byte("b"))

	assert.IsType(t, recordstore.StorageOperation{}, op)

}

func TestMustNewRecordStore(t *testing.T) {
	db := testtools.NewTempBoltDB()
	defer db.Close()
	recordStore := recordstore.MustNewRecordStore(db)
	assert.IsType(t, &recordstore.RecordStore{}, recordStore)
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(recordstore.ResoRecordBucketName))
		assert.NotNil(t, bucket)
		return nil
	})
}

func TestConvertRecordToKeyValues(t *testing.T) {
	resoRecord := recordstore.ResoRecord{
		Identifier:    "a",
		AltIdentifier: []string{"b", "c"},
		OldIdentifier: []string{},
	}

	operations, err := recordstore.ConvertRecordToStorageOperations(resoRecord)
	assert.NoError(t, err)

	assert.Equal(t, recordstore.ResoRecordBucketName, operations[0].Bucket)
	assert.Equal(t, []byte(fmt.Sprintf("%s:a", recordstore.IdentifierKeyPrefix)), operations[0].Key)
	assert.Equal(t, "\x82\xaaIdentifier\xa1a\xadAltIdentifier\x92\xa1b\xa1c", string(operations[0].Value))

	assert.Equal(t, recordstore.ResoRecordBucketName, operations[1].Bucket)
	assert.Equal(t, fmt.Sprintf("%s:a:b", recordstore.AltIdentifierKeyPrefix), string(operations[1].Key))
	assert.Equal(t, "a", string(operations[1].Value))

	resoRecord = recordstore.ResoRecord{
		Identifier:    "a",
		AltIdentifier: []string{"b", "c"},
		OldIdentifier: []string{"d", "e"},
	}

	operations, err = recordstore.ConvertRecordToStorageOperations(resoRecord)
	assert.Len(t, operations, 5)
}

func TestHandleOperation(t *testing.T) {
	operation := recordstore.NewStorageOperation(
		recordstore.ResoRecordBucketName,
		"testing",
		"mykey",
		[]byte("myvalue"))

	db := testtools.NewTempBoltDB()
	defer db.Close()
	recordstore.MustNewRecordStore(db)
	db.Update(func(tx *bolt.Tx) error {
		err := recordstore.HandleOperation(tx, operation)
		assert.NoError(t, err)

		return nil
	})

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(recordstore.ResoRecordBucketName))
		value := bucket.Get([]byte("testing:mykey"))
		assert.Equal(t, []byte("myvalue"), value)
		return nil
	})

}

func TestSaveChunk(t *testing.T) {
	db := testtools.NewTempBoltDB()
	defer db.Close()
	rs := recordstore.MustNewRecordStore(db)
	resoRecord := recordstore.ResoRecord{
		Identifier:    "a",
		AltIdentifier: []string{"b", "c"},
		OldIdentifier: []string{},
	}
	err := rs.SaveChunk([]recordstore.ResoRecord{resoRecord})
	assert.NoError(t, err)
	foundResoRecord, err := rs.FindByIdentifier("a")
	assert.NoError(t, err)
	assert.NotNil(t, foundResoRecord)
	assert.Equal(t, resoRecord.Identifier, foundResoRecord.Identifier)
	assert.Equal(t, resoRecord.AltIdentifier[0], foundResoRecord.AltIdentifier[0])
	assert.Equal(t, len(resoRecord.OldIdentifier), len(foundResoRecord.OldIdentifier))
}
