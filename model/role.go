package model

import "github.com/bearllfleed/go_shop/global"

type Role struct {
	global.GSModel
	Name     string
	ParentId uint64
}

func (Role) TableName() string {
	return "role"
}
