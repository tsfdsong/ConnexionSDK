package model

import (
	"errors"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/db/mysql"
	"reflect"
)

func querySome(m interface{}, list interface{}) error {
	r := mysql.GetDB().Where(m).Find(list)
	return r.Error
}

//	m condition, one result
func queryOne(m interface{}, one interface{}) (err error) {
	return mysql.GetDB().Where(m).First(one).Error
}

func queryLastOne(one interface{}) (err error) {
	return mysql.GetDB().Last(one).Error
}

func queryAll(list interface{}) (total int64, err error) {
	db := mysql.GetDB().Find(list)
	db.Count(&total)
	err = db.Error
	return
}

//	update multiple fields
func update(m interface{}, where interface{}) (err error) {
	db := mysql.GetDB().Model(where).Updates(m)
	if err = db.Error; err != nil {
		return
	}
	if db.RowsAffected != 1 {
		return errors.New("id is invalid and resource is not found")
	}
	return nil
}

//	must ensure it’s primary field has value
func delete(m interface{}) (err error) {
	//WARNING When delete a record, you need to ensure it’s primary field has value, and GORM will use the primary key to delete the record, if primary field’s blank, GORM will delete all records for the models
	//primary key must be not zero value
	db := mysql.GetDB().Delete(m)
	if err = db.Error; err != nil {
		return
	}
	if db.RowsAffected != 1 {
		return errors.New("resource is not found to destroy")
	}
	return nil
}

func getType(v interface{}) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	}
	return t.Name()
}

func getSelected(selected []string, m interface{}, one interface{}) error {
	return mysql.GetDB().Select(selected).Where(m).First(one).Error
}
