package app

import (
	"log"
	"os"

	"github.com/alimasry/gopad/internal/database"
	"github.com/alimasry/gopad/internal/pkg/ot"
	"github.com/alimasry/gopad/internal/pkg/ws"
	"github.com/alimasry/gopad/internal/transport/rest"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var Hub *ws.Hub

func Run() {
	// run the websocket hub loop to handle sending
	Hub = ws.GetHubInstance()
	go Hub.Run()

	// run the operational transformation processing loop
	otbm := ot.GetOTBufferManager()
	go otbm.ProcessTransformations()

	// load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	// initialize database and migrate tables
	if err := database.Init(); err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("APP_PORT")

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	rest.AddRoutes(router)

	router.Run(":" + port)
}
