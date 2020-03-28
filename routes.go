package game

import (
	"github.com/gin-gonic/gin"
)

func (client Client) AddRoutes(prefix string, engine *gin.Engine) *gin.Engine {
	g1 := engine.Group(prefix)

	// Index
	g1.GET("/:status", client.index(prefix))

	// JSON Data for Index
	g1.POST("/:status/json", client.JSONIndexAction)

	// Index
	g1.GET("/:status/user/:uid", client.index(prefix))

	// JSON Data for Index
	g1.POST("/:status/user/:uid/json", client.JSONIndexAction)

	g1.GET("/:status/notifications", client.dailyNotifications)

	return engine
}
