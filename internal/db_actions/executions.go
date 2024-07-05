package dbactions

import (
	"errors"

	"github.com/dopeCape/schduler/internal/models"
	rdb "github.com/dopeCape/schduler/pkg/db"
)

func CreateExecution(execution *models.Execution) (*models.Execution, error) {
	db, err := rdb.GetDb()
	if err != nil {
		return nil, errors.New("Failed to connecte to db")
	}
	res := db.Create(execution)
	if res.Error != nil {
		return nil, err
	}
	return execution, nil

}

func UpdateExecution(execution *models.Execution) error {
	db, err := rdb.GetDb()
	if err != nil {
		return errors.New("Failed to connecte to db")
	}
	res := db.Model(execution).Updates(execution)
	if res.Error != nil {
		return err
	}
	return nil
}

func DeleteExecutionsForTask(taskId string) error {
	db, err := rdb.GetDb()
	if err != nil {
		return errors.New("Failed to connecte to db")
	}
	res := db.Where("TASK_ID = ?", taskId).Delete(&models.Execution{})
	if res.Error != nil {
		return err
	}
	return nil
}

func GetExecutionsForTask(taskID string) ([]models.Execution, error) {
	db, err := rdb.GetDb()
	if err != nil {
		return nil, errors.New("Failed to connecte to db")
	}
	var executions []models.Execution
	res := db.Where("TASK_ID = ?", taskID).Find(&executions)
	if res.Error != nil {
		return nil, err
	}
	return executions, nil

}
