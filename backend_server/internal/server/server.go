package server

import (
	"backend_server/configs"
	"backend_server/pkg/connection/postgres"
	"context"
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	cfg *configs.Config
	db  postgres.DBops
	gin *gin.Engine
}

func NewServer(cfg *configs.Config, db postgres.DBops) *Server {
	return &Server{
		cfg: cfg,
		db:  db,
		gin: gin.Default(),
	}
}
func (s *Server) Run() error {

	s.gin.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))

	err := s.MapHandlers()
	if err != nil {
		return err
	}
	srv := &http.Server{
		Addr:    s.cfg.Server.Host,
		Handler: s.gin,
	}
	go func() {
		s.gin.GET("/health_check", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
		log.Printf("Server is started ", s.cfg.Server.Host)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")

	return nil
}
