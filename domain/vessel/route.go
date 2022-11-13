package vessel

import (
	"kapalapi/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func (d *VesselDeps) VesselRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())
	api.Use(middleware.IsHeader)
	// api.Get("/vessel/v2", d.FindByUserKey)
	// api.Post("/vessel/v2", d.Create)
	api.Post("/vessel", d.CreateNew)
	api.Use(middleware.SecretKey)
	api.Get("/vessel", d.GetVessel)
}
