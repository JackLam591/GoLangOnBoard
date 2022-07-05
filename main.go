package main

import (
	"fmt"
	"helloworld/configs"
	"helloworld/middleware"
	"helloworld/routes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello World!!!!")

	mapErr := godotenv.Load(".env")
	if mapErr != nil {
		fmt.Println(mapErr)
	}

	// start fiber
	app := fiber.New()
	configs.ConnectDB()

	middleware.SetMiddleware(app)
	routes.UserRoute(app)

	app.Listen(":" + os.Getenv("GO_PORT"))
}
