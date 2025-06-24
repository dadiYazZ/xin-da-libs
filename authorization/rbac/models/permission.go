package models

import (
	"errors"
	"github.com/dadiYazZ/xin-da-libs/database"
	"github.com/dadiYazZ/xin-da-libs/object"
	"github.com/dadiYazZ/xin-da-libs/security"
	"gorm.io/gorm"
)

// TableName overrides the table name used by Permission to `profiles`
func (mdl *Permission) TableName() string {
	return mdl.GetTableName(true)
}

// Permission 数据表结构
type Permission struct {
	*database.PowerCompactModel

	PermissionModule *PermissionModule `gorm:"ForeignKey:ModuleID;references:UniqueID" json:"permissionModule"`

	UniqueID    string  `gorm:"column:index_permission_id;index:,unique" json:"permissionID"`
	ObjectAlias *string `gorm:"column:object_alias" json:"objectAlias"`
	ObjectValue string  `gorm:"column:object_value; not null;" json:"objectValue"`
	Action      string  `gorm:"column:action; not null;" json:"action"`
	Description *string `gorm:"column:description" json:"description"`
	ModuleID    *string `gorm:"column:module_id" json:"moduleID"`
}

const TABLE_NAME_PERMISSION = "rbac_permissions"

const PERMISSION_UNIQUE_ID = "index_permission_id"

var TABLE_FULL_NAME_PERMISSION string = "public.ac_" + TABLE_NAME_PERMISSION

const PERMISSION_TYPE_NORMAL int8 = 1
const PERMISSION_TYPE_MODULE int8 = 2

func NewPermission(mapObject *object.Collection) *Permission {

	if mapObject == nil {
		mapObject = object.NewCollection(&object.HashMap{})
	}

	newPermission := &Permission{
		PowerCompactModel: database.NewPowerCompactModel(),
		ObjectAlias:       mapObject.GetStringPointer("objectAlias", ""),
		ObjectValue:       mapObject.GetString("objectValue", ""),
		Action:            mapObject.GetString("action", ""),
		Description:       mapObject.GetStringPointer("description", ""),
		ModuleID:          mapObject.GetStringPointer("moduleID", ""),
	}
	newPermission.UniqueID = newPermission.GetComposedUniqueID()

	return newPermission

}

// 获取当前 Model 的数据库表名称
func (mdl *Permission) GetTableName(needFull bool) string {
	if needFull {
		return TABLE_FULL_NAME_PERMISSION
	} else {
		return TABLE_NAME_PERMISSION
	}
}
func (mdl *Permission) SetTableFullName(tableName string) {
	TABLE_FULL_NAME_PERMISSION = tableName
}

func (mdl *Permission) GetForeignKey() string {
	return "index_permission_id"
}

func (mdl *Permission) GetForeignValue() string {
	return mdl.UniqueID
}

func (mdl *Permission) GetComposedUniqueID() string {

	strKey := mdl.Action + "-" + mdl.ObjectValue
	//fmt2.Dump(strKey)
	hashKey := security.HashStringData(strKey)

	return hashKey
}

func (mdl *Permission) GetRBACRuleName() string {

	return *mdl.ObjectAlias + "-" + mdl.UniqueID[0:5]

}

func (mdl *Permission) CheckPermissionNameAvailable(db *gorm.DB) (err error) {

	result := db.
		//Debug().
		Where("object_alias", mdl.ObjectAlias).
		Where("index_permission_id != ?", mdl.UniqueID).
		First(&Permission{})

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if result.Error != nil {
		return result.Error
	}

	err = errors.New("permission name is not available")

	return err
}
