package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wafi11/high-avaibility/config"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) RegisterRoutes(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		connStr, _ := config.GetConnectionPrimary()

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":       "Documentasi High Availability",
			"status":      "HEALTHY",
			"leader_name": "pg-node-1",
			"primary":     connStr,
			"raw_json":    `{"conn_url": "...", "state": "running"}`,
		})
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}
