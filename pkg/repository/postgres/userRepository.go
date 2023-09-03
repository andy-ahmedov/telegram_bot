package postgres

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/andy-ahmedov/telegram_bot/pkg/logging"
	_ "github.com/lib/pq"
)

type UserRepository struct {
	userDB *sql.DB
	logger *logging.Logger
}

func NewUserRepository(userDB *sql.DB, logger *logging.Logger) *UserRepository {
	return &UserRepository{
		userDB: userDB,
		logger: logger,
	}
}

func (u *UserRepository) GetNumber(chatID int64) (string, error) {
	request := "SELECT phone_number FROM user_repository WHERE chatid = $1"

	row := u.userDB.QueryRow(request, int(chatID))
	var number string
	err := row.Scan(&number)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		u.logger.Error("Ошибка в функции row.Scan ", err)
		return "", err
	}
	return number, nil
}

func (u *UserRepository) GetBonus(chatID int64, phoneNumber string) (int, error) {
	return 0, nil
}

func (u *UserRepository) Save(chatID int64, fieldName string, data interface{}) error {
	query := "SELECT chatid FROM user_repository WHERE chatid = $1"
	row := u.userDB.QueryRow(query, chatID)

	var chat int
	err := row.Scan(&chat)
	if err == sql.ErrNoRows {
		err := u.InsertQuery(chatID, fieldName, data)

		if err != nil {
			return err
		}

	} else if err != nil {
		u.logger.Error("Ошибка в функции row.Scan ", err)
		return err

	} else {
		err := u.UpdateQuery(chatID, fieldName, data)

		if err != nil {
			return err
		}
	}

	return nil
}

func (u *UserRepository) UpdateQuery(chatID int64, field_name string, data interface{}) error {
	reflType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)

	var dataType string
	if reflType.Name() == "string" {
		dataType = "text"
	} else {
		dataType = reflType.Name()
	}

	request := "UPDATE user_repository SET %s = $1::%s WHERE chatid = $2"
	query := fmt.Sprintf(request, field_name, dataType)

	stmt, err := u.userDB.Prepare(query)
	if err != nil {
		u.logger.Error("Ошибка в функции u.userDB.Prepare ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(dataValue.Interface(), int(chatID))
	if err != nil {
		u.logger.Error("Ошибка в функции stmt.Exec ", err)
		return err
	}

	return nil
}

func (u *UserRepository) InsertQuery(chatID int64, field_name string, data interface{}) error {
	reflType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)

	var dataType string
	if reflType.Name() == "string" {
		dataType = "text"
	} else {
		dataType = reflType.Name()
	}

	request := "INSERT INTO user_repository (chatid, %s) VALUES ($1, $2::%s)"
	query := fmt.Sprintf(request, field_name, dataType)

	stmt, err := u.userDB.Prepare(query)
	if err != nil {
		u.logger.Error("Ошибка в функции u.userDB.Prepare ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(int(chatID), dataValue.Interface())
	if err != nil {
		u.logger.Error("Ошибка в функции stmt.Exec ", err)
		return err
	}

	return nil
}

func (u *UserRepository) GetCode(chatID int64) (string, error) {
	request := "SELECT code FROM user_repository WHERE chatid = $1"

	row := u.userDB.QueryRow(request, int(chatID))
	var code string
	err := row.Scan(&code)
	if err != nil {
		if err == sql.ErrNoRows {
			u.logger.Info("Код не найден")
			return "", fmt.Errorf("Code not found")
		}
		return "", err
	}
	return code, nil
}

func (u *UserRepository) ExistCodeInUserRep(chatID int64) (bool, error) {
	query := "SELECT chatid, code FROM user_repository WHERE chatid = $1"
	row := u.userDB.QueryRow(query, int(chatID))

	var chat int
	var code sql.NullString

	err := row.Scan(&chat, &code)
	if err != nil {
		if err == sql.ErrNoRows {
			u.logger.Info("Строка с заданным chatID не найдена")
			return false, nil
		}
		u.logger.Error("Ошибка в функции row.Scan ", err)
		return false, err
	}

	if code.Valid {
		u.logger.Info("Код обнаружен в базе пользователей")
		return true, nil
	}

	return false, nil
}

func (u *UserRepository) DeleteInfo(chatID int64, fieldName string) (bool, error) {
	query := fmt.Sprintf("SELECT chatid, %s FROM user_repository WHERE chatid = $1", fieldName)
	row := u.userDB.QueryRow(query, int(chatID))

	var chat int
	var field sql.NullString
	err := row.Scan(&chat, &field)
	if err != nil {
		if err == sql.ErrNoRows {
			u.logger.Info("Строка с заданным chatID не найдена")
			return true, nil
		}
		return false, err
	}

	if field.Valid {
		var nullstr sql.NullString
		nullstr.Valid = false
		nullstr.String = ""

		// Выполняем SQL-запрос для обновления поля
		// request := fmt.Sprintf("UPDATE user_repository SET %s = $1 WHERE chatid = $2", fieldName)
		// _, err = u.userDB.Exec(request, nullstr, chatID)
		_, err = u.userDB.Exec("DELETE FROM user_repository WHERE chatid = $1", chatID)
		if err != nil {
			u.logger.Error("Ошибка в функции u.userDB.Exec ", err)
			return false, err
		}
	}

	return true, nil
}
