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
	User      *user.Client
	afterLoad bool
}

func NewClient(dsClient *datastore.Client, userClient *user.Client, logger *log.Logger, mcache *cache.Cache, router *gin.Engine, prefix string, afterLoad bool) *Client {
	cl := &Client{
		Client:    sn.NewClient(dsClient, logger, mcache, router),
		User:      userClient,
		afterLoad: afterLoad,
	}
	return cl.addRoutes(prefix)
}

func (cl *Client) addRoutes(prefix string) *Client {
	g1 := cl.Router.Group(prefix)

	// Index
	g1.GET("/:status",
		// gtype.SetTypes(),
		cl.Index(prefix),
	)

	// JSON Data for Index
	g1.POST("/:status/json",
		// gtype.SetTypes(),
		cl.GetFiltered(gtype.All),
		cl.JSONIndexAction,
	)

	// JSON Data for Index
	g1.GET("/:status/json",
		cl.JIndex,
	)

	// Index
	g1.GET("/:status/user/:uid",
		// gtype.SetTypes(),
		cl.Index(prefix),
	)

	// JSON Data for Index
	g1.POST("/:status/user/:uid/json",
		// gtype.SetTypes(),
		cl.GetFiltered(gtype.All),
		cl.JSONIndexAction,
	)

	g1.GET("/:status/notifications",
		cl.GetRunning,
		cl.DailyNotifications,
	)
	return cl
}
