package main

import (
	"github.com/IljaN/chat/user"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Firewall(m *user.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Chat-Auth")
		isAuthenticated, err := m.Authenticated(token)

		if err != nil {
			log.Print(err.Error())
			c.AbortWithError(http.StatusInternalServerError, err)

		}

		if !isAuthenticated {
			c.Abort()
			return
		}

	}

}
