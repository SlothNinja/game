package game

import (
	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/log"
	"github.com/SlothNinja/sn"
	gtype "github.com/SlothNinja/type"
	"github.com/SlothNinja/user"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

type Client struct {
	*sn.Client
	User *user.Client
}

func NewClient(dsClient *datastore.Client, userClient *user.Client, logger *log.Logger, mcache *cache.Cache, router *gin.Engine, prefix string) *Client {
	cl := &Client{
		Client: sn.NewClient(dsClient, logger, mcache, router),
		User:   userClient,
	}
	return cl.addRoutes(prefix)
}

func (client *Client) addRoutes(prefix string) *Client {
	g1 := client.Router.Group(prefix)

	// Index
	g1.GET("/:status",
		// gtype.SetTypes(),
		client.Index(prefix),
	)

	// JSON Data for Index
	g1.POST("/:status/json",
		// gtype.SetTypes(),
		client.GetFiltered(gtype.All),
		client.JSONIndexAction,
	)

	// Index
	g1.GET("/:status/user/:uid",
		// gtype.SetTypes(),
		client.Index(prefix),
	)

	// JSON Data for Index
	g1.POST("/:status/user/:uid/json",
		// gtype.SetTypes(),
		client.GetFiltered(gtype.All),
		client.JSONIndexAction,
	)

	g1.GET("/:status/notifications",
		client.GetRunning,
		client.DailyNotifications,
	)
	return client
}
