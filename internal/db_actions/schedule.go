package dbactions

import (
	"errors"

	"github.com/dopeCape/schduler/internal/models"
	rdb "github.com/dopeCape/schduler/pkg/db"
)

func GetScheduleFromApiKey(prefix string, limit int, offset int) ([]models.Task, error) {
	db, err := rdb.GetDb()
	if err != nil {
		return nil, errors.New("Failed to connecte to db")
	}
	var user models.User
	var tasks []models.Task
	res := db.Where("PREFIX = ?", prefix).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	res = db.Where("USER_ID = ? AND IS_CRON = ?", user.ID, true).Limit(limit).Offset(offset).Find(&tasks)
	if res.Error != nil {
		return nil, res.Error
	}

	return tasks, nil

}
