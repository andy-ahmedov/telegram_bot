package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	commandStart     = "start"
	commandDelete    = "delete"
	commandDiscounts = "discounts"
	commandStores    = "stores"
	commandBalance   = "balance"
	commandNovelties = "novelties"
	commandCatalog   = "catalog"
	commandContacts  = "contacts"
)

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	token, err := b.getAccessToken(message.Chat.ID)
	if token != "" {
		return b.sendText(message.Chat.ID, b.cfg.Messages.StartWhenAuth)
	}

	if err != nil {
		return b.initErrorHandling(errGetDataFromDB, err)
	}

	return b.sendText(message.Chat.ID, b.cfg.Messages.TextStart)
}

func (b *Bot) handleDeleteNumberCommand(message *tgbotapi.Message) error {
	token, _ := b.getAccessToken(message.Chat.ID)
	if token != "" {
		return b.userIsAlreadyLoggedIn(message.Chat.ID)
	}

	_, err := b.userStorage.DeleteInfo(message.Chat.ID, "phone_number")
	if err != nil {
		return b.initErrorHandling(errDeleteDataFromDB, err)
	}

	_, err = b.userStorage.DeleteInfo(message.Chat.ID, "code")
	if err != nil {
		return b.initErrorHandling(errDeleteDataFromDB, err)
	}

	err = b.sendText(message.Chat.ID, b.cfg.Messages.DeleteSuccess)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleDiscountsCommand(message *tgbotapi.Message) error {
	response := ""
	discounts, err := parseDiscounts(b.cfg.Discounts_url)
	if err != nil {
		return b.initErrorHandling(errInternalServer, err)
	}

	for i, discount := range discounts {
		response = response + fmt.Sprintf("%d)	", i+1) + discount + b.cfg.Delimiter
	}

	return b.sendText(message.Chat.ID, response)
}

func (b *Bot) handleBalanceCommand(message *tgbotapi.Message) error {
	token, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return b.initErrorHandling(errGetDataFromDB, err)
	}

	if token == "" {
		return b.sendText(message.Chat.ID, b.cfg.Messages.Unauthorized)
	}

	number, err := b.userStorage.GetNumber(message.Chat.ID)
	if err != nil {
		return b.initErrorHandling(errGetDataFromDB, err)
	}

	balanceResponse, err := b.ComposeBalanceResponse(message.Chat.ID, number)
	if err != nil {
		return err
	}

	return b.sendText(message.Chat.ID, balanceResponse)
}

func (b *Bot) handleStoresCommand(message *tgbotapi.Message) error {
	keyboard := b.menuForStoreCommand(message)

	err := b.handleMenuButton(message, keyboard)
	if err != nil {
		return errSendMsg
	}
	return nil
}

func (b *Bot) handleNoveltiesCommand(message *tgbotapi.Message) error {
	return b.sendText(message.Chat.ID, b.cfg.Messages.Novelties)
}

func (b *Bot) handleCatalogCommand(message *tgbotapi.Message) error {
	pdfFile := b.cfg.Catalog

	file := tgbotapi.NewDocumentUpload(message.Chat.ID, pdfFile)

	_, err := b.bot.Send(file)
	if err != nil {
		b.initErrorHandling(errSendMsg, err)
	}

	return nil
}

func (b *Bot) handleContactsCommand(message *tgbotapi.Message) error {
	return b.sendText(message.Chat.ID, b.cfg.Messages.Contacts)
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	_, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return b.userIsNotAuthorized()
	}

	err = b.sendText(message.Chat.ID, "Неизвестная команда")
	if err != nil {
		return err
	}
	return nil
}
