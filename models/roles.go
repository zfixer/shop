package models

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	"shop/validates"
)

type Role struct {
	gorm.Model

	Name        string `gorm:"unique;not null VARCHAR(191)"`
	DisplayName string `gorm:"VARCHAR(191)"`
	Description string `gorm:"VARCHAR(191)"`
}

/**
 * 通过 id 获取 role 记录
 * @method GetRoleById
 * @param  {[type]}       role  *Role [description]
 */
func GetRoleById(id uint) *Role {
	role := new(Role)
	IsNotFound(DB.Where("id = ?", id).First(role).Error)
	return role
}

/**
 * 通过 name 获取 role 记录
 * @method GetRoleByName
 * @param  {[type]}       role  *Role [description]
 */
func GetRoleByName(name string) *Role {
	role := new(Role)
	IsNotFound(DB.Where("name = ?", name).First(role).Error)
	return role
}

/**
 * 通过 id 删除角色
 * @method DeleteRoleById
 */
func DeleteRoleById(id uint) {
	u := &Role{
		Model: gorm.Model{
			ID: id,
		},
	}
	u.ID = id
	if err := DB.Delete(u).Error; err != nil {
		color.Red(fmt.Sprintf("DeleteRoleErr:%s \n", err))
	}
}

/**
 * 获取所有的角色
 * @method GetAllRole
 * @param  {[type]} name string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAllRoles(name, orderBy string, offset, limit int) (roles []*Role) {

	if err := GetAll(name, orderBy, offset, limit).Find(&roles).Error; err != nil {
		color.Red(fmt.Sprintf("GetAllRoleErr:%s \n", err))
	}
	return
}

/**
 * 创建
 * @method CreateRole
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func CreateRole(aul *validates.RoleRequest, permIds []uint) (role *Role) {
	role = &Role{
		Name:        aul.Name,
		DisplayName: aul.DisplayName,
		Description: aul.Description,
	}

	if err := DB.Create(role).Error; err != nil {
		color.Red(fmt.Sprintf("CreateRoleErr:%v \n", err))
	}

	addPerms(permIds, role)

	return
}

func addPerms(permIds []uint, role *Role) {
	if len(permIds) > 0 {
		roleId := strconv.FormatUint(uint64(role.ID), 10)
		if _, err := Enforcer.DeletePermissionsForUser(roleId); err != nil {
			color.Red(fmt.Sprintf("AppendPermsErr:%s \n", err))
		}
		var perms []Permission
		// 根据制定的权限ID集合， 找到所有权限
		DB.Where("id in (?)", permIds).Find(&perms)
		for _, perm := range perms {
			// 在casbin数据库，为指定角色添加权限集，addpolicy实际作用就是addPermission
			if _, err := Enforcer.AddPolicy(roleId, perm.Name, perm.Act); err != nil {
				color.Red(fmt.Sprintf("AddPolicy:%s \n", err))
			}
		}
	}
}

/**
 * 更新
 * @method UpdateRole
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func UpdateRole(rj *validates.RoleRequest, id uint, permIds []uint) (role *Role) {
	role = &Role{
		Model: gorm.Model{
			ID: id,
		},
	}

	if err := DB.Model(&role).Updates(rj).Error; err != nil {
		color.Red(fmt.Sprintf("UpdatRoleErr:%s \n", err))
	}

	addPerms(permIds, role)

	return
}

// 角色权限
func RolePermisions(id uint) []*Permission {
	perms := Enforcer.GetPermissionsForUser(strconv.FormatUint(uint64(id), 10))
	var ps []*Permission
	for _, perm := range perms {
		if len(perm) >= 3 && len(perm[1]) > 0 && len(perm[2]) > 0 {
			p := GetPermissionByNameAct(perm[1], perm[2])
			if p.ID > 0 {
				ps = append(ps, p)
			}
		}
	}
	return ps
}


