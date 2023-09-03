package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			err := b.handleCommand(update.Message)
			if err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}
			continue
		}

		b.handleMessage(update.Message)
	}
}

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	b.logs.Infof("[%s]: %s", message.From.UserName, message.Text)

	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)

	case commandDelete:
		return b.handleDeleteNumberCommand(message)

	case commandDiscounts:
		return b.handleDiscountsCommand(message)

	case commandStores:
		return b.handleStoresCommand(message)

	case commandBalance:
		return b.handleBalanceCommand(message)

	case commandNovelties:
		return b.handleNoveltiesCommand(message)

	case commandCatalog:
		return b.handleCatalogCommand(message)

	case commandContacts:
		return b.handleContactsCommand(message)

	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	b.logs.Infof("[%s]: %s", message.From.UserName, message.Text)

	if thisCode(message.Text) || thisNumber(message.Text) {
		_, err := b.initAuthorizationProcess(message)
		if err != nil {
			b.handleError(message.Chat.ID, err)
			return
		}
	} else {
		err := b.IsReservedWords(message)
		if err == nil {
			return
		} else if err == errSendMsg {
			b.handleError(message.Chat.ID, errSendMsg)
		}
	}

}

func (b *Bot) userIsAlreadyLoggedIn(chatID int64) error {
	return b.sendText(chatID, b.cfg.Messages.AlreadyLogged)
}

func (b *Bot) userIsNotAuthorized() error {
	return nil
}
