package main

import (
	"ScootinAboot/api"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	api.AddScooterHandlerEndpoints(r)

	r.Run(":8000")
}
