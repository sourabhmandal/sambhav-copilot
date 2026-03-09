package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"sambhavhr/internal/general"
	"sambhavhr/internal/repository"
	"sambhavhr/internal/user"
	"sambhavhr/pkg/database"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"

	env "github.com/caarlos0/env/v11"
)

type EnvConfig struct {
	DBDatabase               string `env:"DB_DATABASE"`
	DBUsername               string `env:"DB_USERNAME"`
	DBPassword               string `env:"DB_PASSWORD"`
	DBHost                   string `env:"DB_HOST"`
	DBPort                   int    `env:"DB_PORT"`
	DBSchema                 string `env:"DB_SCHEMA"`
	ServerPort               int    `env:"SERVER_PORT"`
	GoogleTranslateAPIKey    string `env:"GOOGLE_TRANSLATE_API_KEY"`
	GoogleTranslateProjectID string `env:"GOOGLE_TRANSLATE_PROJECT_ID"`
}

var ENV EnvConfig = EnvConfig{}

func LoadEnvConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file:", err)
	}

	if err := env.ParseWithOptions(&ENV, env.Options{
		RequiredIfNoDef: true,
	}); err != nil {
		fmt.Println(err)
	}
}

func main() {
	// Env config loading
	LoadEnvConfig()
	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	dbInst, db := database.NewDatabasePg(ENV.DBUsername, ENV.DBPassword, ENV.DBHost, ENV.DBDatabase, ENV.DBSchema, ENV.DBPort)

	newServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", ENV.ServerPort),
		Handler: registerRoutes(dbInst, db),
	}
	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(done, newServer, dbInst)

	// start the server
	if err := newServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}

func registerRoutes(dbInst database.Database, db *pgx.Conn) *gin.Engine {
	// Declare Router
	queries := repository.New(db)

	// declare generic handlers
	generalHandlers := general.NewGeneralHandler(dbInst)
	// declare user handlers
	userService := user.NewUserService(queries)
	userHandlers := user.NewUserHandler(userService)

	router := gin.Default()

	// Allow CORS for all origins
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	// generic routes
	router.GET("/health", generalHandlers.HealthCheck)
	// user routes
	userRouter := router.Group("/user")
	userRouter.POST("/", userHandlers.RegisterUser)
	userRouter.GET("/", userHandlers.GetAllUsers)

	return router
}

func gracefulShutdown(done chan bool, server *http.Server, db database.Database) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// shut down any database connections
	if err := db.Close(); err != nil {
		log.Printf("Database unable to stop with error: %v", err)
	}

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}
