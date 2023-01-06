package main

import (
	"bytes"
	"context"
	"database/sql"
	"kapalapi/domain/vessel"
	"kapalapi/pkg"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html"
	"github.com/jackc/pgx/v4"

	_ "github.com/lib/pq"
)

func main() {

	// then := time.Now().AddDate(0, 0, -1).Format("2022-01-01")
	// two := time.Now().Add(-24 * time.Hour).Format("2022-01-01")
	// fmt.Println(then)
	// fmt.Println(two)
	// fmt.Println(time.Now().AddDate(0, 0, -1).Format("2006-01-02"))

	// var arr []vessel.UserKey
	// var i int64
	// for i = 0; i < 4; i++ {
	// 	var val vessel.UserKey
	// 	val.Id = i
	// 	val.Uuid = fmt.Sprintf("row ke %d", i)
	// 	arr = append(arr, val)
	// }
	// fmt.Println(arr)

	// ctx := context.Background()
	ctx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer shutdownCancel()

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "3000"
	}

	databaseUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		databaseUrl = "postgres://root@localhost:5432/kapalapi?sslmode=disable"
	}

	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}

	conn, err := pgx.Connect(ctx, databaseUrl)
	if err != nil {
		log.Fatalf("error connecting to database server : %v", err)
	}

	errMigrate := pkg.Migration(ctx, db)
	if errMigrate != nil {
		log.Fatal(errMigrate)
	}

	app := fiber.New()
	//set cors
	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: "*",
	// 	// AllowOrigins: "https://wahidinaji.github.io/fetch-api-with-github-page/,https://kapalapi.takakiyo.my.id,https://takakiyo.my.id,http://kapalapi.takakiyo.my.id,http://takakiyo.my.id,https://kapalapi-cakrawala.koyeb.app/vessel-keys",
	// 	AllowHeaders: "Origin, Content-Type, Accept,Bearer",
	// 	AllowMethods: "GET, POST, PATCH, PUT, DELETE",
	// }))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://wahidinaji.github.io/fetch-api-with-github-page/,https://kapalapi.takakiyo.my.id,https://takakiyo.my.id,http://kapalapi.takakiyo.my.id,http://takakiyo.my.id,https://kapalapi-cakrawala.koyeb.app/vessel-keys",
		AllowMethods:     "GET, POST, OPTIONS, PUT, DELETE",
		AllowCredentials: true,
		MaxAge:           2592000,
	}))

	app = fiber.New(fiber.Config{
		Views: html.New("./domain/vessel", ".html"),
	})
	app.Get("/vessel-keys", func(c *fiber.Ctx) (err error) {
		c.Render("index", fiber.Map{})
		return
	})

	vessel := vessel.VesselDeps{DB: conn, PQ: db}
	vessel.VesselRoutes(app)

	app.Get("/upload/:", func(c *fiber.Ctx) error {
		fullUri := c.Request().URI().FullURI()
		c.Accepts("image/png")
		c.Accepts("png")
		fullUriString := bytes.NewBuffer(fullUri).String()
		getPathUri := strings.Split(fullUriString, c.BaseURL())

		return c.SendFile("." + getPathUri[1])
	})

	app.Listen(":" + port)
}
