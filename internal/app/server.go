package app

import (
	"net/http"
	"time"

	"github.com/arieffadhlan/go-fitbyte/internal/config"
	"github.com/arieffadhlan/go-fitbyte/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func NewServer(cfg *config.Config, db *sqlx.DB) *http.Server {
	switch cfg.App.Env {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()
	handlers.SetupRouter(cfg, db, router)

	return &http.Server{
		Addr:         cfg.App.Port,
		Handler:      router,
		WriteTimeout: time.Second * 600,
		ReadTimeout:  time.Second * 600,
		IdleTimeout:  time.Second * 600,
	}
}
