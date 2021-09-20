package model

import (
	"next-terminal/pkg/proxy"
	"next-terminal/server/utils"
)

type Proxy struct {
	ID       string         `gorm:"primary_key" json:"id"`
	Name     string         `gorm:"type:varchar(191)" json:"name"`
	Type     proxy.Type     `gorm:"type:varchar(191)" json:"type"`
	Host     string         `gorm:"type:varchar(191)" json:"host"`
	Port     int            `gorm:"type:SMALLINT UNSIGNED" json:"port"`
	Username string         `gorm:"type:varchar(191)" json:"username"`
	Password string         `gorm:"type:varchar(191)" json:"password"`
	Created  utils.JsonTime `json:"created"`
}

func (r *Proxy) TableName() string {
	return "proxies"
}
