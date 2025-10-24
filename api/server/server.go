package server

import (
	"fmt"
	"net/http"
	"os"
	"log/slog"

	"github.com/gin-gonic/gin"
	"panda.com/api/service"
)

type Server struct {
	port   int
	engine *gin.Engine
}

func NewServer(port int, service *service.Service) *Server {
	server := &Server{
		port:   port,
		engine: gin.Default(),
	}

	server.registerRoutes(service)

	return server
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}

func (s *Server) registerRoutes(service *service.Service) {
	s.engine.Use(handleError(),cors())
	
	g := s.engine.Group("/api")

	g.GET("ping", s.handle(service.Ping))
}

func (s *Server) Run() {
	fmt.Println("server running on port", s.port)
	if err := s.engine.Run(fmt.Sprintf(":%d", s.port)); err != nil {
		slog.Error("server run error", err)
		os.Exit(1)
	}
}