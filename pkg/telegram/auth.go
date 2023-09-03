package telegram

import (
	_ "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"time"

	"github.com/andy-ahmedov/telegram_bot/pkg/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) initAuthorizationProcess(message *tgbotapi.Message) (string, error) {

	result := b.checkUserStatus(message)
	err := errors.New("Default error")
	switch result {

	case "expectedCode":
		err = b.needVerificationCode(message.Chat.ID)

	case "verificationCodes":
		err = b.checkVerificationSmsCode(message.Chat.ID, message)

	case "needPhoneNumber":
		err = b.needPhoneNumber(message.Chat.ID)

	case "phoneNumberFound":
		err = b.phoneVerificationProcess(message.Chat.ID, message)

	case "phoneNumberNotFound":
		err = b.accessIsDenied(message.Chat.ID)

	default:
		return "", b.initErrorHandling(errCheckUserStatus, err)
	}

	if err != nil {
		return "", b.handleError(message.Chat.ID, err)
	}

	return "success", nil
}

func (b *Bot) checkUserStatus(message *tgbotapi.Message) string {
	exist, err := b.userStorage.ExistCodeInUserRep(message.Chat.ID)
	if err != nil {
		return fmt.Sprint(err)
	}
	thisNum := thisNumber(message.Text)

	if exist && thisNum {
		return "expectedCode"
	} else if exist && !thisNum {
		return "verificationCodes"
	}

	if !thisNum {
		return "needPhoneNumber"
	}

	exist, err = b.clientStorage.ExistInClientRep(message.Text)
	if err != nil {
		return fmt.Sprint(err)
	}

	if exist {
		return "phoneNumberFound"
	} else {
		return "phoneNumberNotFound"
	}
}

func (b *Bot) phoneVerificationProcess(chatID int64, message *tgbotapi.Message) error {
	verificationCode := generateVerificationCode()
	b.logs.Info("Код подвтерждения для номера ", message.Text, " -> ", verificationCode)

	// if message.Text == "79991946655" { // ДЛЯ ЭКОНОМИИ ДЕНЕГ
	// 	b.userStorage.Save(chatID, "phone_number", message.Text) // ДЛЯ ЭКОНОМИИ ДЕНЕГ
	// 	b.userStorage.Save(chatID, "code", "0000")               // ДЛЯ ЭКОНОМИИ ДЕНЕГ
	// 	b.sendText(chatID, b.config.Messages.CodeSent)           // ДЛЯ ЭКОНОМИИ ДЕНЕГ
	// 	return nil                                               // ДЛЯ ЭКОНОМИИ ДЕНЕГ
	// } // ДЛЯ ЭКОНОМИИ ДЕНЕГ

	if token, _ := b.tokenStorage.Get(message.Text, repository.AccessTokens); token != "" {
		return errNumberAlreadyAuth
	}

	err := b.sendSMStoPhoneNumber(message.Text, verificationCode)
	if err != nil {
		return b.initErrorHandling(errSendSMStoNumber, err)
	}

	err = b.userStorage.Save(chatID, "phone_number", message.Text)
	if err != nil {
		return b.initErrorHandling(errSavingToDB, err)
	}

	err = b.userStorage.Save(chatID, "code", verificationCode)
	if err != nil {
		return b.initErrorHandling(errSavingToDB, err)
	}

	return b.sendText(chatID, b.cfg.Messages.CodeSent)
}

func (b *Bot) userAuthProcess(chatID int64) error {
	token, err := b.createAccessToken()

	if err != nil {
		return b.initErrorHandling(errCreatToken, err)
	}

	phoneNumber, err := b.userStorage.GetNumber(chatID)
	if err != nil {
		return b.initErrorHandling(errGetDataFromDB, err)
	}

	err = b.tokenStorage.Save(phoneNumber, token, repository.AccessTokens)
	if err != nil {
		return b.initErrorHandling(errSavingToDB, err)
	}

	err = b.userStorage.Save(chatID, "authorized", "yes")
	if err != nil {
		return b.initErrorHandling(errSavingToDB, err)
	}

	return b.sendText(chatID, b.cfg.Messages.SuccessAuth)
}

func (b *Bot) checkVerificationSmsCode(chatID int64, message *tgbotapi.Message) error {
	sentCode, err := b.userStorage.GetCode(chatID)
	if err != nil {
		return b.initErrorHandling(errGetDataFromDB, err)
	}

	receivedCode := message.Text

	if sentCode == receivedCode {
		return b.userAuthProcess(chatID)
	}

	return errIncorrectCode
}

func (b *Bot) sendSMStoPhoneNumber(phoneNumber string, code string) error {
	client := &http.Client{}

	requestLink := fmt.Sprintf(b.cfg.SMSlink, b.cfg.SMSapi, phoneNumber, code)
	req, err := http.NewRequest("GET", requestLink, nil)
	if err != nil {
		b.logs.Error("Ошибка в функции http.NewRequest ", err)
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		b.logs.Error("Ошибка в функции client.Do ", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		b.logs.Error("Ошибка в функции ioutil.ReadAll ", err)
		return err
	}
	b.logs.Info(string(body))

	return nil
}

func (b *Bot) needPhoneNumber(chatID int64) error {
	return b.sendText(chatID, b.cfg.Messages.NeedNumber)
}

func (b *Bot) needVerificationCode(chatID int64) error {
	return b.sendText(chatID, b.cfg.Messages.NeedCode)
}

func (b *Bot) accessIsDenied(chatID int64) error {
	return b.sendText(chatID, b.cfg.Messages.AccessDenied)
}

func (b *Bot) getAccessToken(chatID int64) (string, error) {
	phoneNumber, err := b.userStorage.GetNumber(chatID)
	if err != nil {
		return "", b.initErrorHandling(errGetDataFromDB, err)
	}

	if phoneNumber == "" {
		b.logs.Info("Номер не найден")
		return "", nil
	}

	return b.tokenStorage.Get(phoneNumber, repository.AccessTokens)
}

func (b *Bot) createAccessToken() (string, error) {
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)

	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(randomBytes)

	hashString := hex.EncodeToString(hash[:])

	b.logs.Info("Сгенерирован уникальный хеш-ключ")

	return hashString, nil
}

func generateVerificationCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := r.Intn(9999)

	codeStr := fmt.Sprintf("%04d", code)

	return codeStr
}

func thisCode(message string) bool {
	regex := regexp.MustCompile(`^\d{4}$`)
	return regex.MatchString(message)
}

func thisNumber(text string) bool {
	regex := regexp.MustCompile(`^79\d{9}$`)
	return regex.MatchString(text)
}
