package dbactions

import (
	"errors"
	"strings"

	"github.com/dopeCape/schduler/internal/models"
	rdb "github.com/dopeCape/schduler/pkg/db"
)

func GetUserFromAPIKey(key string) (models.User, error) {
	var user models.User

	keySlice := strings.Split(key, ".")
	db, err := rdb.GetDb()

	if err != nil {
		return user, errors.New("Failed to connecte to db")
	}

	res := db.Where("PREFIX = ?", keySlice[0]).First(&user)
	if res.Error != nil {
		return user, err
	}
	return user, nil
}
