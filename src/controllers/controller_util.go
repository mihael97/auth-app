package controllers

import (
	"log"

	"github.com/gin-gonic/gin"
	"gitlab.com/mihael97/Go-utility/src/web/routes"
)

func InitializeRoutes(engine *gin.Engine) {
	controllers := []routes.RoutesController{
		GetProxyController(),
		GetLoginController(),
	}

	log.Print("Adding controller routes")
	routes.AddControllerRoutesWithFilter(false, engine, controllers)
}
