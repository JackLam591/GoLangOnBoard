package routes

import (
	"helloworld/controllers"

	"github.com/gofiber/fiber/v2"
)

func status(c *fiber.Ctx) error {
	return c.SendString("Server is running, status: OK!")
}
func UserRoute(app *fiber.App) {
	app.Get("/", status)

	app.Get("/api/users/all", controllers.GetAllUser)
	app.Get("/api/users/count", controllers.CountUser)
	app.Get("/api/users/:userId", controllers.GetUser)

	app.Post("/api/users", controllers.CreateUser)
	app.Put("/api/users/:userId", controllers.UpdateUser)

	app.Delete("/api/users/:userId", controllers.DeleteUser)
}
