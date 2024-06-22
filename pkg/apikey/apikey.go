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

func (ak *ApiKeyService) GenerateKey() (*models.ApiKey, error) {
	bytes := make([]byte, 32) // 32 bytes = 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return &models.ApiKey{}, err
	}
	encodedKey := hex.EncodeToString(bytes)
	hashedKey := sha256.Sum256([]byte(encodedKey))
	prefix, err := GenerateUniquePrefix()
	if err != nil {
		return &models.ApiKey{}, err
	}
	db, err := rdb.GetDb()
	if err != nil {
		return &models.ApiKey{}, errors.New("Failed to connect to db")
	}
	uniqueId := xid.New()
	token := fmt.Sprint(prefix, ".", hex.EncodeToString(hashedKey[:]))
	apiToken := &models.ApiKey{Key: token, Prefix: prefix, ID: uniqueId.String()}
	createdToken := db.Create(&apiToken)
	if createdToken.Error != nil {
		return &models.ApiKey{}, createdToken.Error
	}
	return apiToken, nil
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
		apikey := &models.ApiKey{}
		tx := db.Where("prefix = ?", prefix).First(&apikey)
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
