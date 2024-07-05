package models

type User struct {
	ID     string `gorm:"primarykey;unique"`
	Email  string `gorm:"primarykey;unique"`
	Prefix string
	Key    string
	Tasks  []Task
}
