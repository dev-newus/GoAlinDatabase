package main

import (
	"github.com/dev-newus/GoAlinDatabase/src/Database"
	"github.com/dev-newus/GoAlinDatabase/src/Type"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	var DbConfig = Type.Config{
		Host:     "localhost",
		Port:     "3306",
		User:     "root",
		Password: "",
		Database: "alin",
	}
	Database.Connect(&DbConfig)
	app.Listen(":8000")
}
