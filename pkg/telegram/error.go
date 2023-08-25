package telegram

import (
	"errors"
	"fmt"
	"log"
)

var (
	// errNoToken           = errors.New("token not found")                        // Токен не найден
	// errXMLFail           = errors.New("Failed to open XML file")                // Ошибка при открытии XML файла
	errUnauthorized      = errors.New("user is not authorized")                 // Пользователь не авторизован
	errDBConnect         = errors.New("Error connecting to db")                 // Ошибка при подключении к базе данных
	errDecodeXML         = errors.New("Cannot decode xml object")               // Ошибка при декодировании XML файла
	errPhoneNotFound     = errors.New("Phone not found")                        // Телефон не найден
	errPrepareDB         = errors.New("Prepare err")                            // Ошибка при подготовке базы данных
	errExecDB            = errors.New("Exec err")                               // Ошибка при исполнении функции в sql
	errCodeNotFound      = errors.New("Code not found")                         // Код подтверждения не найден
	errInternalServer    = errors.New("InternalServerError")                    // Внутренняя ошибка сервера
	errNumberAlreadyAuth = errors.New("Number already is authorized")           // Этот номер уже зарегистрирован
	errIncorrectCode     = errors.New("Incorrect code")                         // Неверный код
	errUnknown           = errors.New("Error description missing")              // Неизвестная ошибка
	errSendMsg           = errors.New("Error sending message to user")          // Ошибка при отправке сообщения пользователю
	errDeleteDataFromDB  = errors.New("Error while deleting data from DB")      // Ошибка при удалении данных с базы данных
	errCheckUserStatus   = errors.New("Error while cheking user status")        // Ошибка при проверке состояния пользователя
	errSendSMStoNumber   = errors.New("Error when sending SMS to user number")  // Ошбика при отправке смс на номер телефона пользователя
	errSavingToDB        = errors.New("Error saving data to database")          // Ошбика при сохранении данных в базу данных
	errCreatToken        = errors.New("Error creating token")                   // Ошбика при создании токена
	errGetDataFromDB     = errors.New("Error while getting data from database") // Ошбика при получении данных с базы данных
)

func (b *Bot) handleError(chatID int64, err error) error {
	var errSend error

	switch err {
	case errUnauthorized:
		errSend = b.sendText(chatID, b.cfg.Messages.Unauthorized)

	case errCreatToken:
		errSend = b.sendText(chatID, b.cfg.Messages.CantCreateToken)

	case errDecodeXML:
		fmt.Println()

	case errNumberAlreadyAuth:
		errSend = b.sendText(chatID, b.cfg.Messages.NumberAlreadyAuth)

	case errDeleteDataFromDB:
		errSend = b.sendText(chatID, b.cfg.Messages.CantDeleteDataFromDB)

	case errCheckUserStatus:
		errSend = b.sendText(chatID, b.cfg.Messages.CheckUserStatus)

	case errSendSMStoNumber:
		errSend = b.sendText(chatID, b.cfg.Messages.CantSendSMSToPhone)

	case errSavingToDB:
		errSend = b.sendText(chatID, b.cfg.Messages.CantSaveToDB)

	case errGetDataFromDB:
		errSend = b.sendText(chatID, b.cfg.Messages.CantGetDataFromDB)

	case errSendMsg:
		log.Printf(b.cfg.Messages.CantSendMessage)
		return nil

	case nil:
		return nil

	default:
		b.sendText(chatID, b.cfg.Messages.UnknownError)
	}

	b.handleError(chatID, errSend)

	return nil
}

func (b *Bot) initErrorHandling(name error, err error) error {
	log.Println(err)

	return name
}
