package casbinx

import (
	"errors"
	"fmt"
	"github.com/PengShaw/go-common/ginx/ginx-auth/models"
	"github.com/PengShaw/go-common/logger"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"strconv"
)

const modelConf = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
`

var (
	enforcer *casbin.Enforcer
)

func Setup() error {
	a := &Adapter{}
	m, err := model.NewModelFromString(modelConf)
	if err != nil {
		return err
	}

	enforcer, err = casbin.NewEnforcer(m, a)
	enforcer.EnableAutoSave(false)
	return err
}

func CheckPermission(groupID uint, uri, method string) bool {
	ok, err := enforcer.Enforce(strconv.Itoa(int(groupID)), uri, method)
	if err != nil {
		logger.Errorf("check permission failed: %s", err.Error())
		return false
	}
	return ok
}

func AddPermission(groupID, sourceID uint) error {
	err := models.AddPermission(groupID, sourceID)
	if err != nil {
		return err
	}
	p, err := models.GetPermission(groupID, sourceID)
	if err != nil {
		return err
	}
	rule := []string{strconv.Itoa(int(groupID)), p.AuthResource.Uri, p.AuthResource.Method}
	_, err = enforcer.AddPolicy(rule)
	if err != nil {
		_ = models.DeletePermission(groupID, sourceID)
	}
	return err
}

func DeletePermission(groupID, sourceID uint) error {
	p, err := models.GetPermission(groupID, sourceID)
	if err != nil {
		return err
	}
	rule := []string{strconv.Itoa(int(groupID)), p.AuthResource.Uri, p.AuthResource.Method}
	_, err = enforcer.RemovePolicy(rule)
	if err != nil {
		return err
	}
	return models.DeletePermission(groupID, sourceID)
}

type Adapter struct{}

func (sa *Adapter) LoadPolicy(model model.Model) error {
	permissions, err := models.ListPermission()
	if err != nil {
		return err
	}
	for _, p := range permissions {
		line := fmt.Sprintf("p,%d,%s,%s", p.AuthUserGroupID, p.AuthResource.Uri, p.AuthResource.Method)
		persist.LoadPolicyLine(line, model)
	}
	return nil
}

func (sa *Adapter) SavePolicy(model model.Model) error {
	// only focus p model
	for _, ast := range model["p"] {
		for _, rule := range ast.Policy {
			groupID, resourceID, err := parseRule(rule)
			if err != nil {
				return err
			}
			err = models.AddPermission(groupID, resourceID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (sa *Adapter) AddPolicy(_ string, _ string, _ []string) error {
	return errors.New("not implemented")
}

func (sa *Adapter) RemovePolicy(_ string, _ string, _ []string) error {
	return errors.New("not implemented")
}

func (sa *Adapter) RemoveFilteredPolicy(_ string, _ string, _ int, _ ...string) error {
	return errors.New("not implemented")
}

func parseRule(rule []string) (groupID, resourceID uint, err error) {
	if len(rule) != 3 {
		return 0, 0, errors.New("len of rule is not 3")
	}
	resource, err := models.GetAuthResource(rule[1], rule[2])
	if err != nil {
		return 0, 0, err
	}
	gid, err := strconv.Atoi(rule[0])
	if err != nil {
		return 0, 0, err
	}
	groupID = uint(gid)
	return groupID, resource.ID, nil
}
