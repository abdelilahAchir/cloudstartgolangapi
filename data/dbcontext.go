package data

import (
	"fmt"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func openMySql(server, database, username, password string, port int) *gorm.DB {
	var url string
	url = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, server, port, database)

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

var Stats = func() (int, int) {
	var antal, wins int64
	DB.Model(&Game{}).Count(&antal)
	DB.Model(&Game{}).Where("winner=?", "You").Count(&wins)
	return int(antal), int(wins)
}

var SaveGame = func(yourSelection, mySelection, winner string) {
	DB.Create(&Game{Winner: winner, YourSelection: yourSelection, MySelection: mySelection, CreatedAt: time.Now()})
}

var InitDatabase = func(file, server, database, username, password string, port int) {
	if len(file) == 0 {
		DB = openMySql(server, database, username, password, port)
	} else {
		DB, _ = gorm.Open(sqlite.Open(file), &gorm.Config{})
	}
	DB.AutoMigrate(&Game{})
}
