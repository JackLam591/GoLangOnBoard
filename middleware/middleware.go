package middleware

import (
	"helloworld/responses"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/timeout"
)

func SetMiddleware(app *fiber.App) {
	app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"jacklam": "1234567",
		},
		Realm: "Forbidden",
		// Authorizer: func(user, pass string) bool {
		// 	if user == "abcd" && pass == "1234567" {
		// 		return true
		// 	}
		// 	return false
		// },
		Unauthorized: func(c *fiber.Ctx) error {
			return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{Status: http.StatusUnauthorized, Message: "Basic Authen Fail", Data: &fiber.Map{"Message": "User Login Fail!"}})
		},
		ContextUsername: "_user",
		ContextPassword: "_pass",
	}))

	app.Use(recover.New())
	app.Get("/", func(c *fiber.Ctx) error {
		panic("Unable to access")
	})

	handler := func(ctx *fiber.Ctx) error {
		err := ctx.SendString("Hello, World 👋!")
		if err != nil {
			return err
		}
		return nil
	}

	app.Get("/foo", timeout.New(handler, 5*time.Second))
}
