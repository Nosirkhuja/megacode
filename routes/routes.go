package routes

import (
	"MegaCode/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", controllers.User)
	app.Post("/logout", controllers.Logout)
	app.Post("/order", controllers.Order)
	app.Get("/oredercollects", controllers.OrderCollects)
}
