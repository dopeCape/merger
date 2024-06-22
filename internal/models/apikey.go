package models

type ApiKey struct {
	ID     string `gorm:"primarykey;unique"`
	Prefix string
	Key    string
}
