package main

import (
	"MegaCode/database"
	logging "MegaCode/internal/pkg/logger"
	"MegaCode/internal/pkg/middleware"
	"MegaCode/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*type config struct {
	HTTPAddr string `envconfig:"HTTP_ADDR" required:"true"`
	Pg       pg.PgConfig
}*/

func main() {
	logger := logging.NewLogger()
	logger.Info().Msgf("starting %v", time.Now())
	err := database.Connect()
	if err != nil {
		logger.Error().Msgf("database error %v", err)
	}
	Middleware := middleware.New("megaCode")
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	app.Use(Middleware.Middleware)
	routes.Setup(app)
	go func() {
		err = app.Listen(":8000")
		if err != nil {
			logger.Error().Msgf("listen error %v", err)
		}
	}()

	srvMetric := http.Server{
		Addr: ":9000",
	}
	go func() {
		r := http.NewServeMux()
		r.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
		srvMetric.Handler = r
		logger.Info().Msgf("metrics start at port %s", srvMetric.Addr)
		if err := srvMetric.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error().Msgf("start metrics server error: %v", err)
			os.Exit(1)
		}
	}()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	)

	for {
		select {
		case <-signals:
			logger.Info().Msg("gracefully shutdown...")
			if err := app.Shutdown(); err != nil {
				logger.Error().Msgf("shutdown error %v", err)
				defer os.Exit(1)
			} else {
				logger.Info().Msg("gracefully stopped")
				return
			}
		}
	}
}
