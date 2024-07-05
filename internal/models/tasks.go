package models

type Status string

const (
	Success  Status = "success"
	Failed   Status = "failed"
	Pending  Status = "pending"
	Archived Status = "archived"
	Active   Status = "active"
)

type Task struct {
	ID            string `gorm:"primarykey;unique"`
	Payload       string
	Headers       []string `gorm:"serializer:json"`
	URL           string
	Queue         string
	Retried       int
	LastErr       string
	Next          string
	CompletedAt   string
	LastErrAt     string
	Status        Status
	SuccessLog    string
	IsCron        bool
	CronExpresion string
	Executions    []Execution
	UserID        string
}

type Execution struct {
	ID          string `gorm:"primarykey;unique"`
	TaskID      string
	Status      Status
	StatusCode  int
	Error       string
	RanAt       string
	CompletedAt string
	SuccessLog  string
	Task        Task
}
