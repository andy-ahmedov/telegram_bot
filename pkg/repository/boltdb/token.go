package boltdb

import (
	"strconv"

	"github.com/andy-ahmedov/telegram_bot/pkg/logging"
	"github.com/andy-ahmedov/telegram_bot/pkg/repository"
	"github.com/boltdb/bolt"
)

type TokenRepository struct {
	db     *bolt.DB
	logger *logging.Logger
}

func NewTokenRepository(db *bolt.DB, logger *logging.Logger) *TokenRepository {
	return &TokenRepository{
		db:     db,
		logger: logger,
	}
}

func (r *TokenRepository) Save(phoneNumber string, token string, bucket repository.Bucket) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put([]byte(phoneNumber), []byte(token))
	})
}

func (r *TokenRepository) Get(phoneNumber string, bucket repository.Bucket) (string, error) {
	var token string

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data := b.Get([]byte(phoneNumber))
		token = string(data)
		return nil
	})

	if err != nil {
		r.logger.Error("Ошибка в функции r.db.View:", err)
		return "", err
	}

	if token == "" {
		r.logger.Info("Токен не найден")
		return "", nil
	}

	return token, nil
}

func intToBytes(value int64) []byte {
	return []byte(strconv.FormatInt(value, 10))
}
