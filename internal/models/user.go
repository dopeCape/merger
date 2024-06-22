package models

type User struct {
	ID   string `gorm:"primarykey;unique"`
	Name string
}
