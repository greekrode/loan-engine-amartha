package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/greekrode/loan-engine-amartha/db"
)

func main() {
	db.InitDB()

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.SetTrustedProxies(nil)

	log.Fatal(router.Run(":8080"))
}
