package main

import (
	"keno/internal/api"
	"keno/internal/db"
	"keno/internal/engine"
	"os"

	// Include Swagger docs in the project
	_ "keno/docs"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Initialise the logger
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})

	// Setup the database
	database, err := db.SetupDatabase("keno.db")
	if err != nil {
		panic(err)
	}

	// Setup the game engine
	gameEngine := engine.SetupEngine(database)

	// Run the API and Engine
	go launchAPI(database, gameEngine)
	gameEngine.StartLoop()
}

// @title           			TAB Keno API
// @version         			1.0
// @description     			This is a sample server for TAB Keno API.
// @host            			localhost:8080
func launchAPI(database *gorm.DB, gameEngine *engine.Engine) {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Use(func(ctx *gin.Context) { ctx.Set(db.DbKey, database) })
	r.Use(func(ctx *gin.Context) { ctx.Set(engine.EngineKey, gameEngine) })

	v1 := r.Group("/api/v1")
	{
		// Make sure the user is authenticated
		v1.Use(api.DiscordAuth)

		// Protected API
		v1.POST("/picks", api.PlacePicks)
		v1.GET("/check/:card_id", api.CheckCard)
	}
	r.GET("/api/v1/ws", api.GameStreamer)
	r.GET("/api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}
