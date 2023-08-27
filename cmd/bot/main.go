package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/andy-ahmedov/telegram_bot/pkg/config"
	"github.com/andy-ahmedov/telegram_bot/pkg/repository"
	"github.com/andy-ahmedov/telegram_bot/pkg/repository/boltdb"
	"github.com/andy-ahmedov/telegram_bot/pkg/repository/postgres"
	"github.com/andy-ahmedov/telegram_bot/pkg/telegram"
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false

	boltDB, err := initBoltDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	postgresDB, err := initPostgresDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer postgresDB.Close()

	tokenRepository := boltdb.NewTokenRepository(boltDB)

	clientRepository := postgres.NewClientRepository(postgresDB)

	userRepository := postgres.NewUserRepository(postgresDB)

	telegramBot := telegram.NewBot(bot, tokenRepository, clientRepository, userRepository, *cfg)

	sigInt := make(chan os.Signal, 1)
	signal.Notify(sigInt, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	for {
		select {
		case sig := <-sigInt:
			log.Printf("Stopping service %v...\n", sig)
			os.Exit(0)
		}
	}
}

func initBoltDB(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return db, nil
}

func initPostgresDB(cfg *config.Config) (*sql.DB, error) {
	postgresInfo := fmt.Sprintf(cfg.ConnectDB, cfg.Host, cfg.Port, cfg.UserName, cfg.PasswordDB, cfg.DBname, cfg.Sslmode)
	db, err := sql.Open("postgres", postgresInfo)
	if err != nil {
		return nil, err
	}
	return db, nil
}
