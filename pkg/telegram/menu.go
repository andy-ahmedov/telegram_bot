package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) menuForStoreCommand(message *tgbotapi.Message) *tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Маяковского"),
			tgbotapi.NewKeyboardButton("Рылеева"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Пушкаревское"),
			tgbotapi.NewKeyboardButton("Оптимус"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Современник"),
			tgbotapi.NewKeyboardButton("Промышленная"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Буинская"),
		),
	)
	return &keyboard
}

// func (b *Bot) menuForDiscountCommand(message *tgbotapi.Message) *tgbotapi.ReplyKeyboardMarkup {
// 	discounts := make([]string, 0)
// 	titleWords := make([]string, 0)
// 	keyboardButton := make([]tgbotapi.KeyboardButton, 0)
// 	doc, err := goquery.NewDocument(url)
// 	if err != nil {
// 		return nil, err
// 	}

// 	doc.Find(".item").Each(func(i int, s *goquery.Selection) {
// 		title := s.Find(".title").Text()
// 		titleWords = append(titleWords, title)
// 		fmt.Println(title)
// 		description := collectText(s.Nodes[0])
// 		discounts = append(discounts, description)
// 	})

// 	return &keyboard
// }

func (b *Bot) handleMenuButton(message *tgbotapi.Message, keyboard *tgbotapi.ReplyKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.cfg.Messages.SelectStore)
	msg.ReplyMarkup = keyboard

	_, err := b.bot.Send(msg)
	if err != nil {
		return b.initErrorHandling(errSendMsg, err)
	}
	return nil
}
