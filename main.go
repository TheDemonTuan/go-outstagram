package main

import (
	"errors"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"os"
	"outstagram/common"
	"outstagram/routes"
)

type Input struct {
	Query         string                 `query:"query"`
	OperationName string                 `query:"operationName"`
	Variables     map[string]interface{} `query:"variables"`
}

func init() {
	common.LoadEnvVar()
	common.ConnectDB()
	common.NewPusherClient()
	//common.CreateStaticFolder(os.Getenv("STATIC_PATH"))
}

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder:       sonic.Marshal,
		JSONDecoder:       sonic.Unmarshal,
		CaseSensitive:     true,
		StrictRouting:     true,
		EnablePrintRoutes: false,
		ServerHeader:      "API Outstagram",
		AppName:           "API Outstagram",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}
			return common.CreateResponse(c,
				code,
				err.Error(),
				nil)
		},
	})

	app.Use(helmet.New())
	app.Use(cors.New(cors.Config{
		//AllowOrigins:  os.Getenv("CLIENT_URL"),
		//ExposeHeaders: os.Getenv("JWT_HEADER"),
	}))
	app.Use(etag.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

	if os.Getenv("APP_ENV") == "development" {
		app.Use(logger.New())
	}

	routes.SetupRouter(app)
	if err := app.Listen(":" + os.Getenv("PORT")); err != nil {
		panic(err)
	}
}
