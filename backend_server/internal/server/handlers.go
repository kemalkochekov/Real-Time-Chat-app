package server

import (
	"backend_server/internal/user/handler/http"
	"backend_server/internal/user/repository/postgres"
	"backend_server/internal/user/usecase"
	"backend_server/internal/ws"
	http_websocket2 "backend_server/internal/ws/handler/http_websocket"
)

func (s *Server) MapHandlers() error {
	userPGRepo := postgres.NewUserRepo(s.cfg, s.db)

	userUC := usecase.NewUserUC(
		s.cfg,
		userPGRepo,
	)
	userHandlers := http.NewUserHandler(s.cfg, userUC)

	userGroup := s.gin.Group("")
	wsGroup := s.gin.Group("ws")
	hub := ws.NewHub()
	go hub.Run()
	wsHandler := http_websocket2.NewHandler(hub)
	http_websocket2.MapUserRoutes(wsGroup, wsHandler)
	http.MapUserRoutes(userGroup, userHandlers)

	return nil
}
