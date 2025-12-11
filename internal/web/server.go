package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/atrox39/ReverseReach/internal/metrics"
	"github.com/atrox39/ReverseReach/internal/tunnel"
)

func Start(port string) {
	r := gin.Default()

	r.Static("/static", "./internal/web/static")

	r.GET("/", func(c *gin.Context) {
		c.File("./internal/web/static/index.html")
	})

	r.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"connections":     metrics.Metrics.Connections,
			"bytes_sent":      metrics.Metrics.BytesSent,
			"bytes_received":  metrics.Metrics.BytesReceived,
			"last_connection": metrics.Metrics.LastConnection,
		})
	})

	r.GET("/logs", func(c *gin.Context) {
		c.JSON(http.StatusOK, tunnel.GetLogEntries())
	})

	r.Run(":" + port)
}
