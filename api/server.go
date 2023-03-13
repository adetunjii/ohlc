package api

import (
	"github.com/adetunjii/ohlc/pkg/logging"
	"github.com/adetunjii/ohlc/service"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	logger logging.Logger
	svc    *service.Service
}

func NewServer(svc *service.Service, logger logging.Logger) *Server {
	server := &Server{
		svc:    svc,
		logger: logger,
	}

	server.setupRouter()
	return server
}

func (s *Server) Start() error {
	return s.router.Run()
}

func (s *Server) setupRouter() {
	router := gin.Default()
	v1 := router.Group("/api/v1")

	s.priceDataRouteGroup(v1)

	s.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
