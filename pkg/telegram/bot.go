package telegram

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/andy-ahmedov/telegram_bot/pkg/config"
	"github.com/andy-ahmedov/telegram_bot/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	bot           *tgbotapi.BotAPI
	tokenStorage  repository.TokenStorage
	clientStorage repository.ClientStorage
	userStorage   repository.UserStorage
	errorStruct   *ErrorStruct
	cfg           config.Config
}

type ErrorStruct struct {
	err  error
	name error
}

func NewBot(bot *tgbotapi.BotAPI, ts repository.TokenStorage, cs repository.ClientStorage, us repository.UserStorage, cfg config.Config) *Bot {
	return &Bot{
		bot:           bot,
		tokenStorage:  ts,
		clientStorage: cs,
		userStorage:   us,
		cfg:           cfg,
	}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account: %s", b.bot.Self.UserName)

	isEmpty, err := b.checkTableIsEmpty("client_repository")
	if isEmpty {
		err = b.clientStorage.DownloadDB(b.cfg.PathToXml)
		if err != nil {
			return b.initErrorHandling(errDecodeXML, err)
		} else {
			log.Println("Database successfully uploaded")
		}
	} else {
		log.Println("Db upload canceled")
	}

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	go b.handleUpdates(updates)
	return nil
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}

func (b *Bot) sendText(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	// msg.ParseMode = tgbotapi.ModeMarkdown
	_, err := b.bot.Send(msg)

	if err != nil {
		return b.initErrorHandling(errSendMsg, err)
	}

	return nil
}

func (b *Bot) checkTableIsEmpty(tableName string) (bool, error) {
	postgresInfo := fmt.Sprintf(b.cfg.ConnectDB, b.cfg.Host, b.cfg.Port, b.cfg.UserName, b.cfg.PasswordDB, b.cfg.DBname, b.cfg.Sslmode)
	db, err := sql.Open("postgres", postgresInfo)
	if err != nil {
		return false, err
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s LIMIT 1);", tableName)
	var exists bool
	err = db.QueryRow(query).Scan(&exists)
	if err != nil {
		return false, err
	}
	return !exists, nil
}
