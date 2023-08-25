package telegram

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/net/html"
)

func collectText(node *html.Node) string {
	var result string

	if node.Type == html.ElementNode && node.Data == "h3" {
		result = getText(node)
		result = "*" + result + "*\n"
	} else if node.Type == html.TextNode {
		result = strings.TrimSpace(node.Data)
		if result != "" {
			result = result + "\n"
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		result += collectText(child)
	}

	return result
}

func getText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	text := ""
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += getText(c)
	}
	return text
}

func extractLinkText(n *html.Node) string {
	var text string

	if n.Type == html.TextNode {
		text = strings.TrimSpace(n.Data)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += extractLinkText(c)
	}

	return text
}

func parseDiscounts(url string) ([]string, error) {
	discounts := make([]string, 0)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	doc.Find(".item").Each(func(i int, s *goquery.Selection) {
		description := collectText(s.Nodes[0])
		discounts = append(discounts, description)
	})

	return discounts, nil
}

func (b *Bot) IsReservedWords(message *tgbotapi.Message) error {
	var errSend error

	switch message.Text {

	case b.cfg.ReservWrds.Mayakovskogo:
		errSend = b.sendText(message.Chat.ID, b.cfg.Messages.Mayakovskogo)

	case b.cfg.ReservWrds.Ryleyeva:
		errSend = b.sendText(message.Chat.ID, b.cfg.Messages.Ryleyeva)

	case b.cfg.ReservWrds.Pushkarovskoye:
		errSend = b.sendText(message.Chat.ID, b.cfg.Messages.Pushkarovskoye)

	case b.cfg.ReservWrds.Optimus:
		errSend = b.sendText(message.Chat.ID, b.cfg.Messages.Optimus)

	case b.cfg.ReservWrds.Sovremennik:
		errSend = b.sendText(message.Chat.ID, b.cfg.Messages.Sovremennik)

	case b.cfg.ReservWrds.Promyshlennaya:
		errSend = b.sendText(message.Chat.ID, b.cfg.Messages.Promyshlennaya)

	case b.cfg.ReservWrds.DiscountCenter:
		errSend = b.sendText(message.Chat.ID, b.cfg.Messages.DiscountCenter)

	default:
		return errors.New("NotReservWords")
	}

	if errSend != nil {
		return errSend
	}

	return nil
}

func (b *Bot) ComposeBalanceResponse(chatID int64, number string) (string, error) {
	balance, err := b.clientStorage.GetBonus(chatID, number)
	if err != nil {
		return "", b.initErrorHandling(errGetDataFromDB, err)
	}

	myBalance := strconv.Itoa(balance)

	fmt.Println(b.cfg.Balance)

	balanceResponse := fmt.Sprintf(b.cfg.Balance, myBalance)

	return balanceResponse, nil
}
