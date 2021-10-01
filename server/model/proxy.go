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
	Username *string        `gorm:"type:varchar(191);default:''" json:"username"`
	Password *string        `gorm:"type:varchar(191);default:''" json:"password"`
	Created  utils.JsonTime `json:"created"`
}

func (p *Proxy) TableName() string {
	return "proxies"
}

func (p *Proxy) ToProxyConfig(dialHost string, dialPort int) *proxy.Config {
	return &proxy.Config{
		Host:     p.Host,
		Port:     p.Port,
		Username: *p.Username,
		Password: *p.Password,
		DialHost: dialHost,
		DialPort: dialPort,
	}
}
