package models

import (
	"github.com/fernandesleticia/go-agenda/database"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type Item struct {
	Id          int `gorm:primary_key`
	Description string
	Done        bool
}

func GetItemByID(Id int) bool {
	item := &Item{}
	result := database.MysqlInstance.First(&item, Id)

	if result.Error != nil {
		log.Warn("Item not found")
		return false
	}

	return true
}

func GetItemsWith(done bool) interface{} {
	var items []Item
	Items := database.MysqlInstance.Where("done = ?", done).Find(&items).Value
	return Items
}
