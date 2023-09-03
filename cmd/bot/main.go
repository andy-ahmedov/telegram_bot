package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/andy-ahmedov/telegram_bot/pkg/config"
	"github.com/andy-ahmedov/telegram_bot/pkg/logging"
	"github.com/andy-ahmedov/telegram_bot/pkg/repository"
	"github.com/andy-ahmedov/telegram_bot/pkg/repository/boltdb"
	"github.com/andy-ahmedov/telegram_bot/pkg/repository/postgres"
	"github.com/andy-ahmedov/telegram_bot/pkg/telegram"
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
)

func main() {
	logger := logging.GetLogger()

	cfg, err := config.Init(logger)
	if err != nil {
		logger.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		logger.Fatal(err)
	}

	bot.Debug = false

	boltDB, err := initBoltDB(cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}

	postgresDB, err := initPostgresDB(cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}
	defer postgresDB.Close()

	tokenRepository := boltdb.NewTokenRepository(boltDB, logger)

	clientRepository := postgres.NewClientRepository(postgresDB, logger)

	userRepository := postgres.NewUserRepository(postgresDB, logger)

	telegramBot := telegram.NewBot(bot, tokenRepository, clientRepository, userRepository, *cfg, logger)

	sigInt := make(chan os.Signal, 1)
	signal.Notify(sigInt, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := telegramBot.Start(); err != nil {
			logger.Fatal(err)
		}
	}()

	for {
		select {
		case sig := <-sigInt:
			logger.Infof("Stopping service %v...\n", sig)
			os.Exit(0)
		}
	}
}

func initBoltDB(cfg *config.Config, logger *logging.Logger) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		logger.Error("Ошибка в функции bolt.Open:", err)
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			logger.Error("Ошибка в функции db.Update:", err)
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}

func initPostgresDB(cfg *config.Config, logger *logging.Logger) (*sql.DB, error) {
	postgresInfo := fmt.Sprintf(cfg.ConnectDB, cfg.Host, cfg.Port, cfg.UserName, cfg.PasswordDB, cfg.DBname, cfg.Sslmode)
	db, err := sql.Open("postgres", postgresInfo)
	if err != nil {
		logger.Error("Ошибка в функции sql.Open:", err)
		return nil, err
	}
	return db, nil
}
