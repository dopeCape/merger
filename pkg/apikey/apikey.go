package apikey

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/dopeCape/schduler/internal/models"
	rdb "github.com/dopeCape/schduler/pkg/db"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type ApiKeyService struct {
}

func NewApiKeySerice() *ApiKeyService {
	return &ApiKeyService{}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (ak *ApiKeyService) GenerateKey(email string) (*models.User, error) {
	bytes := make([]byte, 32) // 32 bytes = 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return &models.User{}, err
	}
	encodedKey := hex.EncodeToString(bytes)
	hashedKey := sha256.Sum256([]byte(encodedKey))
	prefix, err := GenerateUniquePrefix()
	if err != nil {
		return &models.User{}, err
	}
	db, err := rdb.GetDb()
	if err != nil {
		return &models.User{}, errors.New("Failed to connect to db")
	}
	var user models.User
	res := db.Where(&models.User{Email: email}).First(&user)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return &models.User{}, errors.New("Failed to find user")
	}

	token := fmt.Sprint(prefix, ".", encodedKey)
	user.Email = email
	user.Key = hex.EncodeToString(hashedKey[:])
	user.Prefix = prefix
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		uniqueId := xid.New()
		user.ID = uniqueId.String()
		res := db.Create(&user)
		if res.Error != nil {
			return &models.User{}, errors.New("Failed to create user")
		}
	} else {
		res := db.Model(&models.User{}).Where("EMAIL = ?", email).Updates(&user)
		if res.Error != nil {
			fmt.Println(res.Error)
			return &models.User{}, errors.New("Failed to update user")
		}
	}
	user.Key = token
	return &user, nil
}

func GenerateUniquePrefix() (string, error) {
	var prefix string
	db, err := rdb.GetDb()
	if err != nil {
		return "", errors.New("Failed to connect to db")
	}
	for {
		// Generate a random 5-digit number
		prefix, err = randomString(7)
		if err != nil {
			return "", err
		}
		// Check if the prefix is unique in the database
		user := &models.User{}
		tx := db.Where("prefix = ?", prefix).First(&user)
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			break
		}
	}
	return prefix, nil
}

func randomString(length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}
