package storage

import (
	"fmt"
	"log"

	"github.com/casmeyu/micro-user/structs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect(cnf structs.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cnf.Db.User,
		cnf.Db.Password,
		cnf.Db.Ip,
		cnf.Db.Name,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("[storage.database] (Connect) - Error occurred while connecting to database", err.Error())
		return nil, err
	}
	return db, nil
}

func MakeMigration(cnf structs.Config, entity interface{}) error {
	fmt.Println("Making migration from storage funciton")
	db, err := Connect(cnf)
	if err != nil {
		return err
	}

	err = db.AutoMigrate(entity)
	if err != nil {
		log.Println("[storage.database] (MakeMigration) - Error occurred while making a migration", err.Error())
		return err
	}
	log.Println("Migration succesful")
	return nil
}
