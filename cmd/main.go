package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/naufalatha/go-boilerplate/config"
	"github.com/naufalatha/go-boilerplate/database"
	"github.com/naufalatha/go-boilerplate/handlers"
	"github.com/naufalatha/go-boilerplate/helpers/logger"
	"github.com/naufalatha/go-boilerplate/routes"
	"github.com/naufalatha/go-boilerplate/routes/middleware"
	"github.com/rs/zerolog"
)

// GracefulShutdown handle close signal and shutdown the program properly
func GracefulShutdown(app *fiber.App, db *sql.DB, logFile *os.File) chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s
		fmt.Println("Shutting down marketplace-service...")
		app.Shutdown()
		database.Shutdown(db)
		logFile.Close()
		defer func() {
			fmt.Println("Shutdown complete, bye!")
			os.Exit(0)
		}()
		close(wait)
	}()
	return wait
}

// SetupZeroLog initialize zerolog logger
func SetupZeroLog() (*zerolog.Logger, *os.File) {
	file, err := os.OpenFile(".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		panic(err)
	}
	output := zerolog.MultiLevelWriter(os.Stdout, file)
	zeroLog := zerolog.New(output).With().Timestamp().Logger()
	return &zeroLog, file
}

func main() {
	log, file := SetupZeroLog()
	log.Info().Msg("Starting marketplace-service...")
	config := config.LoadConfig(".", log)
	logger.InitLogger(log, config.AppEnv) // Custom logger object to write to file

	// Init fiber
	var app *fiber.App = fiber.New(middleware.SetupFiberConfig(config))
	middleware.Default(app, config, log)

	// Middleware init
	middleware.Default(app, config, log)

	// Init & Inject dependency
	db := database.NewConnection(config)
	database := database.InitDatabase(db)
	handler := handlers.InitHandlers(config, database, log)

	// JWT Guard and Public API Route
	app.Route("", routes.InitRouter(handler, config).Route)    // Public API Route
	middleware.UseJWT(app, config)                             // Inject JWT Middleware
	app.Route("", routes.InitRouter(handler, config).JWTRoute) // JWT Guarded Route

	// Serve HTTP, default port:8080
	wait := GracefulShutdown(app, db, file)
	log.Info().Msg(fmt.Sprintf("Application is running in %s and listening to port %s", config.AppEnv, config.AppPort))
	app.Listen(fmt.Sprintf(":%s", config.AppPort))
	<-wait
}
