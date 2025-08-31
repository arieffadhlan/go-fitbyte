package handlers

import (
	"github.com/arieffadhlan/go-fitbyte/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func SetupRouter(cfg *config.Config, db *sqlx.DB, r *gin.Engine) {

	r.GET("/health-check", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "health 🇮🇩🇮🇩🇮🇩🇮🇩🇮🇩",
		})
	})
}
