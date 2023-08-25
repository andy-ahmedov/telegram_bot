package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	_ "github.com/lib/pq"
)

type UserRepository struct {
	userDB *sql.DB
}

func NewUserRepository(userDB *sql.DB) *UserRepository {
	return &UserRepository{userDB: userDB}
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
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(dataValue.Interface(), int(chatID))
	if err != nil {
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
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(int(chatID), dataValue.Interface())
	if err != nil {
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
			er := fmt.Sprint("Code not found")
			return "", errors.New(er)
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
			// Строка с заданным ID не найдена
			return false, nil
		}
		return false, err
	}

	if code.Valid {
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
			// Строка с заданным ID не найдена
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
			return false, err
		}
	}

	return true, nil
}
