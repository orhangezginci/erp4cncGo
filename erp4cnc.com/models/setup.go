package models
import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
  )

  var DB *gorm.DB

  func ConnectDatabase() {
  	dsn := "host=postgres user=postgres password=postgres dbname=erpcnc port=5432 sslmode=disable TimeZone=Europe/Berlin"
		//dsn := "host=postgres user=postgres password=postgres dbname=teachingminds port=5432 sslmode=disable TimeZone=Europe/Berlin"

	  database, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	database.AutoMigrate(&Machine{})

	DB = database

  }
