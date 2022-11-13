package main

import (
	"bytes"
	"context"
	"kapalapi/domain/vessel"
	"kapalapi/pkg"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5"
)

func main() {

	// var arr []vessel.UserKey
	// var i int64
	// for i = 0; i < 4; i++ {
	// 	var val vessel.UserKey
	// 	val.Id = i
	// 	val.Uuid = fmt.Sprintf("row ke %d", i)
	// 	arr = append(arr, val)
	// }
	// fmt.Println(arr)

	ctx := context.TODO()

	databaseUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		databaseUrl = "postgres://root@localhost:5432/kapalapi?sslmode=disable"
	}

	conn, err := pgx.Connect(ctx, databaseUrl)
	if err != nil {
		log.Fatalf("error connecting to database server : %v", err)
	}

	errMigrate := pkg.Mig(ctx, conn)
	if errMigrate != nil {
		log.Fatal(errMigrate)
	}

	app := fiber.New()
	//set cors
	app.Use(cors.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept,Bearer",
		AllowMethods: "GET, POST, PATCH, PUT, DELETE",
	}))

	vessel := vessel.VesselDeps{DB: conn}
	vessel.VesselRoutes(app)

	app.Get("/upload/:", func(c *fiber.Ctx) error {
		fullUri := c.Request().URI().FullURI()
		c.Accepts("image/png")
		c.Accepts("png")
		fullUriString := bytes.NewBuffer(fullUri).String()
		getPathUri := strings.Split(fullUriString, c.BaseURL())

		return c.SendFile("." + getPathUri[1])
	})

	app.Listen(":3000")
}
