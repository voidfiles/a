package data_manager

import (
	"github.com/asdine/storm"
)

// NodeInterface is a subset of storm.Node
type NodeInterface interface {
	Save(interface{}) error
}

// DataManager will manage database interactions
type DataManager struct {
	db *storm.DB
}

// TransactionFunction will execute inside of a transaction
type TransactionFunction func(NodeInterface) error

// MustNewDataManager will create a new DataManager
func MustNewDataManager(db *storm.DB) *DataManager {
	return &DataManager{
		db: db,
	}
}

// InTransaction will save a list of data in a tx
func (dm *DataManager) InTransaction(update TransactionFunction) error {
	tx, err := dm.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = update(tx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// Save will save an item to a database
func (dm *DataManager) Save(data interface{}) error {
	return dm.db.Save(data)
}
