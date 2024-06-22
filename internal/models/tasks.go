package models

import (
	"database/sql/driver"
	"errors"
	"strings"
)

type Status string

const (
	Success  Status = "success"
	Failed   Status = "failed"
	Pending  Status = "pending"
	Archived Status = "archived"
	Active   Status = "active"
)

type MultiString []string

func (s *MultiString) Scan(src interface{}) error {
	str, ok := src.(string)
	if !ok {
		return errors.New("failed to scan multistring field - source is not a string")
	}
	*s = strings.Split(str, ",")
	return nil
}

func (s MultiString) Value() (driver.Value, error) {
	if s == nil || len(s) == 0 {
		return nil, nil
	}
	return strings.Join(s, ","), nil
}

type Task struct {
	ID          string `gorm:"primarykey;unique"`
	Payload     string
	Headers     []string `gorm:"serializer:json"`
	URL         string
	Queue       string
	Retried     int
	LastErr     string
	Next        string
	CompletedAt string
	LastErrAt   string
	// Status        Status `gorm:"type:enum('success','pending','failed','archived','active')"`
	SuccessLog    string
	IsCron        bool
	CronExpresion string
}

type Execution struct {
	ID string `gorm:"primarykey;unique"`
	// Status      Status `gorm:"type:enum('success','pending','failed','archived','active')"`
	ErrorCode   int
	Error       string
	RanAt       string
	CompletedAt string
	SuccessLog  string
}
