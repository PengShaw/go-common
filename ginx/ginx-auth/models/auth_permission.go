package models

import "gorm.io/gorm"

// AuthPermission Table
type AuthPermission struct {
	gorm.Model
	AuthUserGroupID uint          `gorm:"uniqueIndex:auth_user_group_idx_auth_resource_id"`
	AuthResourceID  uint          `gorm:"uniqueIndex:auth_user_group_idx_auth_resource_id"`
	AuthUserGroup   AuthUserGroup `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AuthResource    AuthResource  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func AddPermission(groupID, sourceID uint) error {
	item := AuthPermission{
		AuthUserGroupID: groupID,
		AuthResourceID:  sourceID,
	}
	existed, err := existRecord(&item, "auth_user_group_id = ? AND auth_resource_id = ?", item.AuthUserGroupID, item.AuthResourceID)
	if err != nil {
		return err
	}
	if existed {
		return nil
	}
	return db.Create(&item).Error
}

func DeletePermission(groupID, sourceID uint) error {
	return db.Unscoped().Where("auth_user_group_id = ? AND auth_resource_id = ?", groupID, sourceID).
		Delete(&AuthPermission{}).Error
}

func GetPermission(groupID, sourceID uint) (AuthPermission, error) {
	var item AuthPermission
	err := db.Joins("AuthUserGroup").
		Joins("AuthResource").
		Where("auth_user_group_id = ? AND auth_resource_id = ?", groupID, sourceID).First(&item).Error
	return item, err
}

func GetPermissionByID(id uint) (AuthPermission, error) {
	var item AuthPermission
	err := db.Joins("AuthUserGroup").
		Joins("AuthResource").
		Where(" \"auth_permissions\".\"id\" = ?", id).First(&item).Error
	return item, err
}

func ListPermission() ([]AuthPermission, error) {
	var items []AuthPermission
	err := db.Joins("AuthUserGroup").
		Joins("AuthResource").
		Find(&items).Error
	return items, err
}
