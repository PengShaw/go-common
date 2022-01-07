package models

import (
	"errors"
	"github.com/PengShaw/go-common/hash"
	"github.com/PengShaw/go-common/logger"
	"gorm.io/gorm"
)

// AuthUser Table
type AuthUser struct {
	gorm.Model
	Username        string `gorm:"unique;uniqueIndex:auth_user_group_idx_username"`
	HashedPassword  string
	AuthUserGroupID uint          `gorm:"uniqueIndex:auth_user_group_idx_username"`
	AuthUserGroup   AuthUserGroup `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func VerifyUser(username, password string) (AuthUser, error) {
	var user AuthUser
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Errorf("query user err: %s", err.Error())
		}
		return user, errors.New("wrong username or password")
	}
	if !hash.VerifyHashedPassword(password, user.HashedPassword) {
		return user, errors.New("wrong username or password")
	}
	return user, nil
}

func CreateUser(username, password string, groupID uint) error {
	hashed, err := hash.PasswordHash(password)
	if err != nil {
		return err
	}
	return db.Create(&AuthUser{
		Username:        username,
		HashedPassword:  hashed,
		AuthUserGroupID: groupID,
	}).Error
}

func ListUser() ([]AuthUser, error) {
	var users []AuthUser
	err := db.Joins("AuthUserGroup").Find(&users).Error
	return users, err
}

func DeleteUser(id uint) error {
	return db.Unscoped().Delete(AuthUser{}, id).Error
}

func UserChangePassword(userID uint, password string) error {
	hashed, err := hash.PasswordHash(password)
	if err != nil {
		return err
	}
	return db.Model(&AuthUser{}).Where("id = ?", userID).Update("hashed_password", hashed).Error
}

func RegroupUser(userID, groupID uint) error {
	return db.Model(&AuthUser{}).Where("id = ?", userID).Update("auth_user_group_id", groupID).Error
}

func GetUser(id uint) (AuthUser, error) {
	var item AuthUser
	err := db.Where("id = ?", id).First(&item).Error
	return item, err
}
