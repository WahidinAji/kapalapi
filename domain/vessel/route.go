package vessel

import (
	"kapalapi/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func (d *VesselDeps) VesselRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())
	api.Get("/vessel-keys", d.GetAll)
	// api.Get("/vessel/date/:from<datetime(2006\\-01\\-02)>/:to<datetime(2006\\-01\\-02)>", d.GetVesselBydate)
	api.Get("/vessel/date", d.GetVesselBydate)
	api.Use(middleware.IsHeader)
	// api.Get("/vessel/v2", d.FindByUserKey)
	// api.Post("/vessel/v2", d.Create)
	api.Post("/vessel", d.CreateNew)
	api.Use(middleware.SecretKey)
	api.Get("/vessel", d.GetVessel)
}
