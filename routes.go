package game

import (
	"cloud.google.com/go/datastore"
	gtype "github.com/SlothNinja/type"
	"github.com/gin-gonic/gin"
)

type Client struct {
	DS *datastore.Client
}

func NewClient(dsClient *datastore.Client) Client {
	return Client{
		DS: dsClient,
	}
}

func (client Client) AddRoutes(prefix string, engine *gin.Engine) *gin.Engine {
	g1 := engine.Group(prefix)

	// Index
	g1.GET("/:status",
		gtype.SetTypes(),
		client.Index(prefix),
	)

	// JSON Data for Index
	g1.POST("/:status/json",
		gtype.SetTypes(),
		client.GetFiltered(gtype.All),
		client.JSONIndexAction,
	)

	// Index
	g1.GET("/:status/user/:uid",
		gtype.SetTypes(),
		client.Index(prefix),
	)

	// JSON Data for Index
	g1.POST("/:status/user/:uid/json",
		gtype.SetTypes(),
		client.GetFiltered(gtype.All),
		client.JSONIndexAction,
	)

	g1.GET("/:status/notifications",
		client.GetRunning,
		client.DailyNotifications,
	)

	return engine
}
