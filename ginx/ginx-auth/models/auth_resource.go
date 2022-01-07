package models

import (
	"gorm.io/gorm"
)

// AuthResource Table
type AuthResource struct {
	gorm.Model
	Uri    string `gorm:"uniqueIndex:uri_x_method"`
	Method string `gorm:"uniqueIndex:uri_x_method"`
}

func setupAuthResource(items []AuthResource) error {
	for _, v := range items {
		existed, err := existRecord(&v, "uri = ? AND method = ?", v.Uri, v.Method)
		if err != nil {
			return err
		}
		if existed {
			continue
		}
		err = db.Create(&v).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func ListAuthResource() ([]AuthResource, error) {
	var items []AuthResource
	err := db.Find(&items).Error
	return items, err
}

func GetAuthResource(uri, method string) (AuthResource, error) {
	var item AuthResource
	err := db.Where("uri = ? AND method = ?", uri, method).First(&item).Error
	return item, err
}
