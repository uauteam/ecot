package repo

import (
	"github.com/jinzhu/gorm"
	"github.com/uauteam/ecot/err"
)

var DBMapping map[string]*gorm.DB

func RegisterDB(dbName string, db *gorm.DB) error {
	if DBMapping == nil {
		DBMapping = make(map[string]*gorm.DB)
	}

	if _, ok := DBMapping[dbName]; ok {
		db = DBMapping[dbName]
		return err.DBAlreadyRegistered
	}

	DBMapping[dbName] = db

	return nil
}

func DB(dbName string) *gorm.DB {
	d, ok := DBMapping[dbName]
	if !ok {
		d = &gorm.DB{}
		d.AddError(err.DBNotRegistered)
		return d
	}

	return d
}
