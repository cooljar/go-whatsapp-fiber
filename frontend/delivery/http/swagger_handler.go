package http

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

func NewSwaggerHandler(app *fiber.App) {
	// Routes for GET method:

	//default
	app.Get("/swagger/*", swagger.Handler) // get one user by ID

	/*app.Get("/swagger/*", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})*/

	//custom
	/*route.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL: "http://example.com/doc.json",
		DeepLinking: false,
	}))*/
}
