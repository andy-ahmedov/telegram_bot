package repository

import "encoding/xml"

type Bucket string

const (
	AccessTokens Bucket = "access_tokens"
)

type CatalogObject struct {
	XMLName  xml.Name    `xml:"CatalogObject.ФизическиеЛица"`
	Surname  string      `xml:"Description"`
	ContInfo ContactInfo `xml:"КонтактнаяИнформация"`
}

type ContactInfo struct {
	PhoneNumber string `xml:"НомерТелефона"`
}

const (
	CatalogObjectName = "CatalogObject.ФизическиеЛица"
)

type TokenStorage interface {
	Get(phoneNumber string, bucket Bucket) (string, error)
	Save(phoneNumber string, token string, bucket Bucket) error
}

type ClientStorage interface {
	DownloadDB(file string, tableName string) error
	FindNumber(phoneNumber string) (bool, error)
	ExistInClientRep(number string) (bool, error)
	CreateTable(tableName string, createTableCode string) error
	GetBonus(chatID int64, phoneNumber string) (int, error)
	UpdateDB(sqlCode string, file string) error
}

type UserStorage interface {
	GetCode(chatID int64) (string, error)
	GetNumber(chatID int64) (string, error)
	ExistCodeInUserRep(chatID int64) (bool, error)
	GetBonus(chatID int64, phoneNumber string) (int, error)
	DeleteInfo(chatID int64, fieldName string) (bool, error)
	Save(chatID int64, fieldName string, data interface{}) error
}
