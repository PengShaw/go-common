package models

import (
	"gorm.io/gorm"
)

// AuthUserGroup Table
type AuthUserGroup struct {
	gorm.Model
	GroupName string `gorm:"unique"`
	Comment   string
}

func CreateUserGroup(name, comment string) error {
	return db.Create(&AuthUserGroup{
		GroupName: name,
		Comment:   comment,
	}).Error
}

func DeleteUserGroup(id uint) error {
	return db.Unscoped().Delete(AuthUserGroup{}, id).Error
}

func ListUserGroup() ([]AuthUserGroup, error) {
	var items []AuthUserGroup
	err := db.Find(&items).Error
	return items, err
}

func RenameUserGroup(id uint, groupName string) error {
	return db.Model(&AuthUserGroup{}).Where("id = ?", id).Update("group_name", groupName).Error
}

func ChangeUserGroupComment(id uint, comment string) error {
	return db.Model(&AuthUserGroup{}).Where("id = ?", id).Update("comment", comment).Error
}

func GetUserGroup(id uint) (AuthUserGroup, error) {
	var item AuthUserGroup
	err := db.Where("id = ?", id).First(&item).Error
	return item, err
}
