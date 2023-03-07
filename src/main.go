package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mihael97/auth-proxy/src/controllers"
	"gitlab.com/mihael97/Go-utility/src/env"
)

func main() {
	port := env.GetEnvVariable("HTTP_PORT", "8080")
	log.Printf("Auth proxy starting at port %s", port)

	router := gin.Default()
	controllers.InitializeRoutes(router)

	log.Panic(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
