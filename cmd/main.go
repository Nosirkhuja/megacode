package main

import (
	"MegaCode/database"
	logging "MegaCode/internal/pkg/logger"
	"MegaCode/internal/pkg/pg"
	"MegaCode/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/kelseyhightower/envconfig"
	"time"
)

type config struct {
	HTTPAddr string `envconfig:"HTTP_ADDR" required:"true"`
	Pg       pg.PgConfig
}

func main() {
	logger := logging.NewLogger()
	logger.Info().Msgf("starting %v", time.Now())
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		logger.Error().Msgf("reading config error %v", err)
	}
	err := database.Connect()
	if err != nil {
		logger.Error().Msgf("database error %v", err)
	}
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	routes.Setup(app)

	err = app.Listen(":8000")
	if err != nil {
		logger.Error().Msgf("listen error %v", err)
	}
}
