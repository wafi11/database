package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/wafi11/high-avaibility/config"
	"github.com/wafi11/high-avaibility/pkg/server"
)

func main() {
	// cfg := config.NewConfig()
	connStr, err := config.GetConnectionPrimary()
	if err != nil {
		fmt.Println("Error getting connection string: ", err)
		return
	}
	db := config.NewDBConfig(connStr)
	db.Connect()

	router := gin.Default()
	server := server.NewServer()
	server.RegisterRoutes(router)

	router.Run(":8080")
}
