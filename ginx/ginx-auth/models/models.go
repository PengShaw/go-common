package models

import (
	"errors"
	contextkey "github.com/PengShaw/go-common/ginx/ginx-context-key"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetupDB(DB *gorm.DB, router *gin.Engine, exceptResource []string) error {
	db = DB

	// setup all table
	if err := db.AutoMigrate(
		&AuthUser{},
		&AuthUserGroup{},
		&AuthResource{},
		&AuthPermission{},
	); err != nil {
		return err
	}

	// setup data of AuthResource Table
	var authResourceList []AuthResource
	for _, v := range router.Routes() {
		skip := false
		for _, v1 := range exceptResource {
			if v1 == v.Path {
				skip = true
			}
		}
		if skip {
			continue
		}
		authResourceList = append(authResourceList, AuthResource{
			Uri:    v.Path,
			Method: v.Method,
		})
	}
	if err := setupAuthResource(authResourceList); err != nil {
		return err
	}

	// setup Admin
	if err := setupAdmin(); err != nil {
		return err
	}
	return nil
}

func existRecord(record interface{}, query interface{}, args ...interface{}) (bool, error) {
	res := db.Where(query, args...).First(record)
	if res.Error == nil {
		return true, nil
	}
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, res.Error
}

func setupAdmin() error {
	var createAdmin = func(createFunc func() error, record interface{}, query interface{}, args ...interface{}) error {
		existed, err := existRecord(record, query, args...)
		if err != nil {
			return err
		}
		if !existed {
			err = createFunc()
			if err != nil {
				return err
			}
			_, err = existRecord(record, query, args...)
			if err != nil {
				return err
			}
		}
		return nil
	}

	// create admin group
	group := AuthUserGroup{}
	err := createAdmin(func() error {
		return CreateUserGroup(contextkey.DefaultAuthAdminUsergroup, contextkey.DefaultAuthAdminUsergroupComment)
	}, &group, "group_name = ?", contextkey.DefaultAuthAdminUsergroup)
	if err != nil {
		return err
	}
	// create admin user
	err = createAdmin(func() error {
		return CreateUser(contextkey.DefaultAuthAdminUsername, contextkey.DefaultAuthAdminPassword, group.ID)
	}, &AuthUser{}, "username = ?", contextkey.DefaultAuthAdminUsername)
	if err != nil {
		return err
	}

	// init admin permission
	resources, err := ListAuthResource()
	if err != nil {
		return err
	}
	for _, v := range resources {
		err := AddPermission(group.ID, v.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
