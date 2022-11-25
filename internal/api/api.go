package api

import (
	"wishmill/internal/config"
	"wishmill/internal/db"
	"wishmill/internal/logger"

	docs "wishmill/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Caltoph API
// @version 1.0
// @BasePath /_caltoph/v1

type errormsg struct {
	Message string `json:"message"`
}

func Init() {
	logger.DebugLogger.Println("api: Initializing api")
	//Initializer gin
	if !config.Config.DevMode {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.GET("/_wishmill/v1/health", getHealth)
	docs.SwaggerInfo.BasePath = "/_wishmill/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	go router.Run(":8080")
	logger.DebugLogger.Println("api: Finished initializing api")
}

// getHealth godoc
// @Summary Get app health
// @Despcription Checks db and cna connection
// @Success 200
// @Failure 500
// @Router /health [get]
func getHealth(c *gin.Context) {
	logger.DebugLogger.Println("api: getHealth is called")
	health := db.GetDbHealth()
	if health {
		c.String(200, "HEALTHY")
		logger.DebugLogger.Println("api: getHealth is OK")
		return
	}
	if !health {
		c.String(500, "NOT HEALTHY")
		logger.WarningLogger.Println("api: getHealth is NOT OK")
		return
	}
}
