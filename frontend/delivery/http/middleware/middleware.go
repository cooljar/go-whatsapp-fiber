package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtMiddleware "github.com/gofiber/jwt/v2"
	"os"
)

// GoMiddleware represent the data-struct for middleware
type GoMiddleware struct {
	appCtx *fiber.App
	// another stuff , may be needed by middleware
}

// CORS will handle the CORS middleware
func (m *GoMiddleware) CORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "*",
		//AllowOrigins: "https://gofiber.io, https://gofiber.net",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, HEAD, PUT, PATCH, POST, DELETE",
	})
}

// LOGGER simple logger.
func (m *GoMiddleware) LOGGER() fiber.Handler {
	return logger.New()
}

// JWT jwt.
func (m *GoMiddleware) JWT() fiber.Handler {
	// Create config for JWT authentication middleware.
	config := jwtMiddleware.Config{
		SigningKey:   []byte(os.Getenv("JWT_SECRET_KEY")),
		ContextKey:   "jwt", // used in private routes
		ErrorHandler: jwtError,
	}

	return jwtMiddleware.New(config)
}

func jwtError(c *fiber.Ctx, err error) error {
	// Return status 400 and failed authentication error.
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 401 and failed authentication error.
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": true,
		"msg":   err.Error(),
	})
}

// InitMiddleware initialize the middleware
func InitMiddleware(ctx *fiber.App) *GoMiddleware {
	return &GoMiddleware{appCtx: ctx}
}
