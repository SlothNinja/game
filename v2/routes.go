package game

import (
	"github.com/SlothNinja/sn/v2"
	"github.com/SlothNinja/user"
)

type Client struct {
	*sn.Client
	User      *user.Client
	afterLoad bool
}

func NewClient(snClient *sn.Client, userClient *user.Client, prefix string, afterLoad bool) *Client {
	cl := &Client{
		Client:    snClient,
		User:      userClient,
		afterLoad: afterLoad,
	}
	return cl.addRoutes(prefix)
}

func (cl *Client) addRoutes(prefix string) *Client {
	// g1 := cl.Router.Group(prefix)

	// // Index
	// g1.GET("/:status",
	// 	// gtype.SetTypes(),
	// 	cl.Index(prefix),
	// )

	// JSON Data for Index
	// g1.POST("", cl.GamesIndex)

	// // JSON Data for Index
	// g1.PUT("/:status/json",
	// 	cl.JIndex,
	// )

	// // Index
	// g1.GET("/:status/user/:uid",
	// 	// gtype.SetTypes(),
	// 	cl.Index(prefix),
	// )

	// // JSON Data for Index
	// g1.POST("/:status/user/:uid/json",
	// 	// gtype.SetTypes(),
	// 	cl.GetFiltered(sn.All),
	// 	cl.JSONIndexAction,
	// )

	// g1.GET("/:status/notifications",
	// 	cl.GetRunning,
	// 	cl.DailyNotifications,
	// )
	return cl
}
