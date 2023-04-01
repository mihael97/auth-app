package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mihael97/auth-proxy/src/controllers"
	config "github.com/mihael97/auth-proxy/src/util"
)

func main() {
	config := config.GetConfig()

	var port string

	if config.Port != nil {
		port = *config.Port
	} else {
		port = "8080"
	}

	log.Printf("Auth proxy starting at port %s", port)

	router := gin.Default()
	controllers.InitializeRoutes(router)

	log.Panic(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
