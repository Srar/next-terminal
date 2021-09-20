package model

import (
	"next-terminal/pkg/proxy"
	"next-terminal/server/utils"
)

type Session struct {
	ID               string         `gorm:"primary_key" json:"id"`
	Protocol         string         `gorm:"type:varchar(191)" json:"protocol"`
	IP               string         `gorm:"type:varchar(191)" json:"ip"`
	Port             int            `gorm:"type:SMALLINT UNSIGNED" json:"port"`
	ProxyType        proxy.Type     `gorm:"type:varchar(191)" json:"proxyType"`
	ProxyHost        string         `gorm:"type:varchar(191)" json:"proxyHost"`
	ProxyPort        int            `gorm:"type:SMALLINT UNSIGNED" json:"proxyPort"`
	ProxyUsername    string         `gorm:"type:varchar(191)"  json:"proxyUsername"`
	ProxyPassword    string         `gorm:"type:varchar(191)"  json:"proxyPassword"`
	ConnectionId     string         `gorm:"type:varchar(191)"  json:"connectionId"`
	AssetId          string         `gorm:"type:varchar(191);index" json:"assetId"`
	Username         string         `gorm:"type:varchar(191)" json:"username"`
	Password         string         `gorm:"type:varchar(191)" json:"password"`
	Creator          string         `gorm:"type:varchar(191);index" json:"creator"`
	ClientIP         string         `gorm:"type:varchar(191)" json:"clientIp"`
	Width            int            `gorm:"type:SMALLINT UNSIGNED" json:"width"`
	Height           int            `gorm:"type:SMALLINT UNSIGNED" json:"height"`
	Status           string         `gorm:"index" json:"status"`
	Recording        string         `json:"recording"`
	PrivateKey       string         `json:"privateKey"`
	Passphrase       string         `json:"passphrase"`
	Code             int            `gorm:"type:SMALLINT" json:"code"`
	Message          string         `json:"message"`
	ConnectedTime    utils.JsonTime `json:"connectedTime"`
	DisconnectedTime utils.JsonTime `json:"disconnectedTime"`
	Mode             string         `gorm:"type:varchar(191)" json:"mode"`
}

func (r *Session) TableName() string {
	return "sessions"
}

type SessionForPage struct {
	ID               string         `json:"id"`
	Protocol         string         `json:"protocol"`
	IP               string         `json:"ip"`
	Port             int            `json:"port"`
	Username         string         `json:"username"`
	ConnectionId     string         `json:"connectionId"`
	AssetId          string         `json:"assetId"`
	Creator          string         `json:"creator"`
	ClientIP         string         `json:"clientIp"`
	Width            int            `json:"width"`
	Height           int            `json:"height"`
	Status           string         `json:"status"`
	Recording        string         `json:"recording"`
	ConnectedTime    utils.JsonTime `json:"connectedTime"`
	DisconnectedTime utils.JsonTime `json:"disconnectedTime"`
	AssetName        string         `json:"assetName"`
	CreatorName      string         `json:"creatorName"`
	Code             int            `json:"code"`
	Message          string         `json:"message"`
	Mode             string         `json:"mode"`
}
