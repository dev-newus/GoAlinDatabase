package Database

import (
	"fmt"
	"time"

	"github.com/dev-newus/GoAlinDatabase/src/Type"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(Config *Type.Config) {
	var err error
	//convert address to string dsn
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		Config.User, Config.Password, Config.Host, Config.Port, Config.Database)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("connect database success")

	//connection pooling

	sqlDB, err := DB.DB()
	if err != nil {
		panic("failed to pooling database")
	}
	sqlDB.SetMaxIdleConns(Config.SetMaxIdleConns)
	sqlDB.SetMaxOpenConns(Config.SetMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(Config.SetConnMaxLifetime) * time.Minute)

	fmt.Println("Max Idle Conns : ", sqlDB.Stats().MaxOpenConnections)
}
