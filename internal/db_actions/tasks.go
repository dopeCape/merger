package dbactions

import (
	"errors"

	"github.com/dopeCape/schduler/internal/models"
	rdb "github.com/dopeCape/schduler/pkg/db"
)

func CreateTask(task *models.Task) (*models.Task, error) {
	db, err := rdb.GetDb()
	if err != nil {
		return nil, errors.New("Failed to connecte to db")
	}
	res := db.Create(task)
	if res.Error != nil {
		return nil, err
	}
	return task, nil

}

func UpdateTask(task *models.Task) error {
	db, err := rdb.GetDb()
	if err != nil {
		return errors.New("Failed to connecte to db")
	}
	res := db.Model(task).Updates(task)
	if res.Error != nil {
		return err
	}
	return nil
}

func DeleteTask(task *models.Task) error {
	db, err := rdb.GetDb()
	if err != nil {
		return errors.New("Failed to connecte to db")
	}
	res := db.Delete(task)
	if res.Error != nil {
		return err
	}
	return nil

}

func GetTaskForAPIKey(prefix string, limit int, offset int) ([]models.Task, error) {
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
	res = db.Where("USER_ID = ?", user.ID).Limit(limit).Offset(offset).Find(&tasks)
	if res.Error != nil {
		return nil, res.Error
	}

	return tasks, nil

}
