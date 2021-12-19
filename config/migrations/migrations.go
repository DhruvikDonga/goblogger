package migrations

import (
	"fmt"

	"github.com/DhruvikDonga/tabsmooth/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

var DNS = "root:@tcp(127.0.0.1:3306)/smarttab?charset=utf8mb4&parseTime=True&loc=Local"

//Add more models in automigrate function to initialize migrations
func IntialMigration() {
	DB, err = gorm.Open(mysql.Open(DNS), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to the DB")
	}
	DB.AutoMigrate(
		&models.Blogs{},
		&models.Users{},
	)

}
