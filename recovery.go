package slogger

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log := slog.Default()
				if l, ok := c.Get("logger"); ok {
					log = l.(*slog.Logger)
				}

				log.Error("panic recovered",
					slog.Any("panic", r),
					slog.String("stack", string(debug.Stack())),
					slog.String("path", c.Request.URL.Path),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "internal server error",
				})
			}
		}()
		c.Next()
	}
}