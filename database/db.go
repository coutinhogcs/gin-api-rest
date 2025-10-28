package database

import (
	"fmt"
	"log"

	"github.com/coutinhogcs/api-go-gin/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func DbConection() {
	strConn := "host=localhost user=admin password=root dbname=alunos port=5433 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(strConn))
	fmt.Println(DB)
	if err != nil {
		log.Panic("Error to access database")
	}

	err := DB.AutoMigrate(&models.Aluno{})
	if err != nil {
		return
	}
}
