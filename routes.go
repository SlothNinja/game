package game

import (
	gtype "github.com/SlothNinja/type"
	"github.com/gin-gonic/gin"
)

func (client Client) AddRoutes(prefix string, engine *gin.Engine) *gin.Engine {
	g1 := engine.Group(prefix)

	// Index
	g1.GET("/:status",
		gtype.SetTypes(),
		client.index(prefix),
	)

	// JSON Data for Index
	g1.POST("/:status/json",
		gtype.SetTypes(),
		client.getFiltered(gtype.All),
		client.jsonIndexAction,
	)

	// Index
	g1.GET("/:status/user/:uid",
		gtype.SetTypes(),
		client.index(prefix),
	)

	// JSON Data for Index
	g1.POST("/:status/user/:uid/json",
		gtype.SetTypes(),
		client.getFiltered(gtype.All),
		client.jsonIndexAction,
	)

	g1.GET("/:status/notifications",
		client.getRunning,
		client.dailyNotifications,
	)

	return engine
}
