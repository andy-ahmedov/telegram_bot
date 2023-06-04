package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	_ "github.com/lib/pq"
)

type myData struct {
	id        int
	name      string
	last_name string
	number    string
	bonus     int
}

type apiBot struct {
	bot      *tgbotapi.BotAPI
	update   tgbotapi.Update
	message  tgbotapi.Message
	chatID   int64
	userName string
}

var dbInfo = fmt.Sprintf("user=andy password=%s dbname=telegram_bot sslmode=disable", PASSWORD)

func collectData(data []string) error {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	if len(data) < 4 {
		return errors.New("Мало данных")
	}
	for i := 0; i < len(data); i++ {
		switch i {
		case 0, 1:
			hasSpecialChars, err := regexp.MatchString("[^A-Za-zА-Яа-я]", data[i])
			if err != nil {
				fmt.Println("Ошибка в MatchString")
				return err
			}
			if hasSpecialChars {
				fmt.Println("Неверный формат ввода имени и фамилии")
				return errors.New("Error")
			}
		case 2:
			// re, err := regexp.Compile(`^7-\d{3}-\d{3}-\d{2}-\d{2}$`)
			matched, err := regexp.MatchString(`^7-\d{3}-\d{3}-\d{2}-\d{2}$`, data[i])
			if err != nil {
				fmt.Println("Ошибка в MatchString")
				return err
			}
			// hasSpecialChars := re.MatchString(data[i])
			if !matched {
				fmt.Println("Неверный формат ввода номера")
				return errors.New("Error")
			}
		case 3:
			_, err := strconv.Atoi(data[i])
			if err != nil {
				fmt.Println("Неверный формат ввода бонусов")
				return err
			}
		default:
			fmt.Println("Много данных в сообщении")
			return errors.New("Error")
		}
	}

	request := "INSERT INTO persons (first_name, last_name, phone_number, bonus) VALUES($1, $2, $3, $4)"

	if _, err = db.Exec(request, data[0], data[1], data[2], data[3]); err != nil {
		fmt.Println("Ошибка при внесении данных в БД")
		return err
	}

	return nil
}

func findData(data string) ([]myData, error) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var dataSlice []myData
	fieldName := []string{"first_name", "last_name", "phone_number"}
	for _, value := range fieldName {
		// request := "SELECT * FROM persons WHERE first_name LIKE 'Эльмаддин%'"
		request := fmt.Sprintf("SELECT * FROM persons WHERE %s ILIKE $1", value)
		rows, err := db.Query(request, data)
		if err != nil {
			log.Println("Ошибка в функции db.Query")
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var slice myData
			err := rows.Scan(&slice.id, &slice.name, &slice.last_name, &slice.number, &slice.bonus)
			if err != nil {
				log.Println("Ошибка в функции rows.Scan")
				return nil, err
			}
			dataSlice = append(dataSlice, slice)
		}
		err = rows.Err()
		if err != nil {
			return nil, err
		}
	}

	return dataSlice, nil
}

func deleteRow(id int) error {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()
	// request := "SELECT * FROM persons WHERE first_name LIKE 'Эльмаддин%'"
	request := fmt.Sprintf("DELETE FROM persons WHERE id = %d", id)
	if _, err = db.Exec(request); err != nil {
		fmt.Println("Ошибка при удалении данных в БД")
		return err
	}

	return nil
}

func showAllData() ([]myData, error) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var dataSlice []myData
	request := "SELECT * FROM persons"
	rows, err := db.Query(request)
	if err != nil {
		log.Println("Ошибка в функции db.Query")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var slice myData
		err := rows.Scan(&slice.id, &slice.name, &slice.last_name, &slice.number, &slice.bonus)
		if err != nil {
			log.Println("Ошибка в функции rows.Scan")
			return nil, err
		}
		dataSlice = append(dataSlice, slice)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return dataSlice, nil
}

func sendText(apiBot *apiBot, text string) {
	msg := tgbotapi.NewMessage(apiBot.chatID, text)
	apiBot.bot.Send(msg)
}

func writeCommand(apiBot *apiBot) error {
	log.Println(apiBot.message.Chat.UserName, "Юзер ввел команду write")
	if len(apiBot.message.Text) < 7 {
		log.Println("Пустая команда write")
		sendText(apiBot, "Введите данные после команды и отправьте снова.")
		return errors.New("Пустая команда")
	}
	result := strings.Split(apiBot.message.Text[7:], ", ")
	log.Println(result)
	if err := collectData(result); err != nil {
		log.Println("Ошибка в функции collectData")
		sendText(apiBot, "Неверный формат ввода данных.")
	} else {
		log.Println("Добавлена запись")
		sendText(apiBot, "Запись успешно добавлена!")

	}
	return nil
}

func deletRowCommand(apiBot *apiBot) error {
	log.Println(apiBot.userName, "Юзер ввел команду delete_row")
	if len(apiBot.message.Text) < 12 {
		log.Println("Пустая команда delete_row")
		sendText(apiBot, "Введите id строки которую хотите удалить.")
		return errors.New("Пустая команда")
	}
	result := apiBot.message.Text[12:]
	log.Println(result)
	id, err := strconv.Atoi(result)
	if err != nil {
		log.Println("Ошибка в функции ATOI")
		sendText(apiBot, "Ошибка на стороне сервера.")
		return errors.New("Ошибка в функции ATOI")
	}
	if err := deleteRow(id); err != nil {
		log.Println("Ошибка в функции deleteRow")
		sendText(apiBot, "Неверный формат ввода")
		log.Fatal(err)
	} else {
		log.Println("Удалена запись", id)
		sendText(apiBot, "Запись успешно удалена!")
	}
	return nil
}

func findCommand(apiBot *apiBot) error {

}

func main() {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	bot, err := tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
		panic(err)
	}

	apiBot := &apiBot{bot: bot}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	log.Printf("Authorized on bot %s", bot.Self.UserName)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if reflect.TypeOf(update.Message.Text).Kind() == reflect.String && update.Message.Text != "" {
			apiBot.chatID = update.Message.Chat.ID
			apiBot.message = *update.Message
			switch update.Message.Command() {
			case "start":
				log.Println(update.Message.Chat.UserName, "Юзер запустил бота")
				sendText(apiBot, "Привет, Бот может выдавать данные по запросу и записывать данные. Пока вы ничего не записывали, база данных пуста.\nЧтобы записать данные необходимо воспользоваться командой '/write'.\nДля получения данных необходимо воспользоваться командой '/find' и написать данные в таком формате:\nИмя(Макс 20 симв), Фамилия(Макс 25 симв), Номер(7-999-999-99-99), Бонус(Число)\nЧтобы показать все записи воспользуйтесь командой '/show_all'")

			case "write":
				if writeCommand(apiBot) != nil {
					continue
				}

			case "delete_row":
				if deletRowCommand(apiBot) != nil {
					continue
				}

			case "find":
				log.Println(update.Message.Chat.UserName, "Юзер ввел команду find")
				trimString := update.Message.Text[6:]
				log.Println(trimString)
				result, err := findData(trimString)
				if err != nil {
					log.Println("Ошибка в функции findData")
					sendText(apiBot, "Ошибка на стороне сервера.")
					log.Fatal(err)
				} else {
					for i := 0; i < len(result); i++ {
						str := fmt.Sprintf("%s %s %s", result[i].name, result[i].last_name, result[i].number)
						sendText(apiBot, str)
					}
				}

			case "show_all":
				log.Println(update.Message.Chat.UserName, "Юзер ввел команду show_all")
				result, err := showAllData()
				if err != nil {
					log.Println("Ошибка в функции showAllData")
					sendText(apiBot, "Ошибка на стороне сервера.")
					log.Fatal(err)
				} else {
					for i := 0; i < len(result); i++ {
						str := fmt.Sprintf("%d) %s %s %s", result[i].id, result[i].name, result[i].last_name, result[i].number)
						sendText(apiBot, str)
					}
				}

			default:
				log.Println(update.Message.Chat.UserName, "Юзер ввел неверную команду")
				sendText(apiBot, "Oups, do you speak עִברִית (Hebrew)?")
			}
		} else {
			log.Println(update.Message.Chat.UserName, "Юзер ввел что-то непонятное")
			sendText(apiBot, "Интересненько, ответить тем же я к сожалению не могу..")
		}
	}
}
