package postgres

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"log"
	"os"

	"github.com/andy-ahmedov/telegram_bot/pkg/repository"
	_ "github.com/lib/pq"
)

type ClientRepository struct {
	clientDB *sql.DB
}

func NewClientRepository(clientDB *sql.DB) *ClientRepository {
	return &ClientRepository{clientDB: clientDB}
}

func (c *ClientRepository) DownloadDB(file string, tableName string) error {
	xmlData, err := os.Open(file)
	if err != nil {
		log.Fatal("Error opening client databasef file: ", err)
	}

	var catalogObject repository.CatalogObject

	d := xml.NewDecoder(xmlData)
	i := 0
	for t, _ := d.Token(); t != nil; t, _ = d.Token() {
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == repository.CatalogObjectName {
				err := d.DecodeElement(&catalogObject, &se)
				if err != nil {
					log.Fatal("Cannot decode xml object : ", err)
				}
				request := fmt.Sprintf("INSERT INTO %s (client_name, phone_number, bonus) VALUES ($1, $2, 0)", tableName)
				_, err = c.clientDB.Exec(request, catalogObject.Surname, catalogObject.ContInfo.PhoneNumber)
				if err != nil {
					log.Fatal(err)
				}
				i++
			}
		}
	}
	return nil
}

func (c *ClientRepository) FindNumber(phoneNumber string) (bool, error) {
	request := "SELECT * FROM client_repository WHERE phone_number ILIKE $1"
	log.Println("Введенный номер: ", phoneNumber)

	id, bonus := 0, 0
	client_name, phone := "", ""
	if err := c.clientDB.QueryRow(request, phoneNumber).Scan(&id, &client_name, &phone, &bonus); err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("phone not found: %v", phoneNumber)
		}
		return false, fmt.Errorf("error: %v", err)
	}
	log.Println(client_name, phone)
	return true, nil
}

func (c *ClientRepository) GetBonus(chatID int64, phoneNumber string) (int, error) {
	if phoneNumber == "79991946655" { // ДЛЯ ТЕСТА
		return 168, nil // ДЛЯ ТЕСТА
	} // ДЛЯ ТЕСТА

	request := "SELECT bonus FROM client_repository WHERE phone_number = $1"

	row := c.clientDB.QueryRow(request, phoneNumber)
	var bonus int
	err := row.Scan(&bonus)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil
		}
		return -1, err
	}
	return bonus, nil
}

func (c *ClientRepository) ExistInClientRep(number string) (bool, error) {
	if number == "79991946655" { // ДЛЯ ТЕСТА
		return true, nil // ДЛЯ ТЕСТА
	} // ДЛЯ ТЕСТА

	query := "SELECT client_name FROM client_repository WHERE phone_number = $1"
	row := c.clientDB.QueryRow(query, number)

	var name string
	err := row.Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (c *ClientRepository) tableExists(tableName string) (bool, error) {
	query := fmt.Sprintf("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = %s);", tableName)
	var exists bool
	err := c.clientDB.QueryRow(query).Scan(&exists)
	if err != nil {
		log.Println(err)
		return false, err
	}
	return exists, nil
}

func (c *ClientRepository) CreateTable(tableName string, createTableCode string) error {
	_, err := c.clientDB.Exec(createTableCode)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (c *ClientRepository) UpdateDB(sqlCode string, file string) error {
	log.Println("Begin rename tables")

	_, err := c.clientDB.Exec(sqlCode)
	if err != nil {
		log.Println(err)
		return err
	}

	err = os.Remove(file)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Client DB successfully updated!")

	_, err = c.clientDB.Exec("DROP TABLE client_repository_old;")
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
