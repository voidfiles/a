package data_manager_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/asdine/storm"
	"github.com/stretchr/testify/assert"
	"github.com/voidfiles/a/data_manager"
)

func createDB(t *testing.T, opts ...func(*storm.Options) error) (*storm.DB, func()) {
	dir, err := ioutil.TempDir(os.TempDir(), "storm")
	if err != nil {
		t.Error(err)
	}
	db, err := storm.Open(filepath.Join(dir, "storm.db"), opts...)
	if err != nil {
		t.Error(err)
	}

	return db, func() {
		db.Close()
		os.RemoveAll(dir)
	}
}

func TestMustNewDataManager(t *testing.T) {

	db, cleanup := createDB(t)
	defer cleanup()

	dm := data_manager.MustNewDataManager(db)
	assert.IsType(t, &data_manager.DataManager{}, dm)
}

type TestStruct struct {
	ID string
}

func TestSave(t *testing.T) {

	db, cleanup := createDB(t)
	defer cleanup()

	dm := data_manager.MustNewDataManager(db)
	testData := &TestStruct{ID: "blah"}
	dm.Save(testData)
}

func TestInTransaction(t *testing.T) {

	db, cleanup := createDB(t)
	defer cleanup()

	dm := data_manager.MustNewDataManager(db)
	dm.InTransaction(func(dbx data_manager.NodeInterface) error {
		testData1 := &TestStruct{ID: "blah1"}
		dbx.Save(testData1)
		testData2 := &TestStruct{ID: "blah2"}
		dbx.Save(testData2)
		return nil
	})
}
