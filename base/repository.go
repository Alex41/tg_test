package base

import (
	"fmt"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var db *gorm.DB

func init() {

	var (
		dbHost = os.Getenv("DB_HOST")
		dbName = os.Getenv("DB_NAME")
		dbUser = os.Getenv("DB_USER")
		dbPass = os.Getenv("DB_PASSWORD")
		dbDSN  = fmt.Sprintf(`host=%s user=%s password=%s sslmode=disable dbname=%s`, dbHost, dbUser, dbPass, dbName)
	)

	var err error
	db, err = gorm.Open(postgres.Open(dbDSN))
	if err != nil {
		panic(err)
	}

	// I used goose migrations because gorm can't migrate enums.
	native, _ := db.DB()
	err = goose.Up(native, "./migrations")
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&Message{})
	if err != nil {
		panic(err)
	}

}

func SaveMessage(message *Message) error {
	return db.Save(message).Error
}

func GetMessages() (m []Message, e error) {
	e = db.Order("id desc").Find(&m).Error
	return
}
