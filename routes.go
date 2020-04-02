package game

import (
	gtype "github.com/SlothNinja/type"
	"github.com/gin-gonic/gin"
)

func AddRoutes(prefix string, engine *gin.Engine) {
	g1 := engine.Group(prefix)

	// Index
	g1.GET("/:status",
		gtype.SetTypes(),
		Index(prefix),
	)

	// JSON Data for Index
	g1.POST("/:status/json",
		gtype.SetTypes(),
		GetFiltered(gtype.All),
		JSONIndexAction,
	)

	// Index
	g1.GET("/:status/user/:uid",
		gtype.SetTypes(),
		Index(prefix),
	)

	// JSON Data for Index
	g1.POST("/:status/user/:uid/json",
		gtype.SetTypes(),
		GetFiltered(gtype.All),
		JSONIndexAction,
	)

	g1.GET("/:status/notifications",
		GetRunning,
		DailyNotifications,
	)
}

//import (
//	"bitbucket.org/SlothNinja/slothninja-games/sn/mlog"
//	"bitbucket.org/SlothNinja/slothninja-games/sn/user"
//
//	"github.com/go-martini/martini"
//)
//
//func AddRoutes(prefix string, r martini.Router) {
//	group := "/" + prefix
//	r.Group(group, func(r martini.Router) {
//		r.Get("/notifications",
//			GetRunning,
//			DailyNotifications,
//		).Name(prefix + "/notifications")
//		addIndexRoute(r, prefix)
//	})
//}
//
////func AddDefaultRoutes(prefix string, r martini.Router, remote ...bool) {
////	remoteServer := false
////	if len(remote) == 1 {
////		remoteServer = remote[0]
////	}
////	if remoteServer {
////		addHomeRoute(r)
////		addUserRoutes(r)
////	}
////	group := "/" + prefix
////	r.Group(group, func(r martini.Router) {
////		prefix2 := "game"
////		group2 := "/" + prefix2
////		r.Group(group2, func(r martini.Router) {
////			prefix := prefix + "_" + prefix2
////			//addNewActionRoute(r, prefix)
////			//addCreateRoute(r, prefix)
////			//addShowRoute(r, prefix)
////			//addUpdateRoute(r, prefix)
////			//addAcceptRoute(r, prefix)
////			//addDropRoute(r, prefix)
////			addAddMessageRoute(r, prefix)
////			//addJSONRoute(r, prefix)
////		})
////		prefix2 = "games"
////		group2 = "/" + prefix2
////		r.Group(group2, func(r martini.Router) {
////			prefix := prefix + "_" + prefix2
////			addIndexRoute(r, prefix)
////		})
////		prefix3 := "admin"
////		group3 := "/" + prefix3
////		r.Group(group3, func(r martini.Router) {
////			prefix := prefix + "_" + prefix3
////			addAdminShowRoute(r, prefix)
////			addAdminUpdateRoute(r, prefix)
////			addAdminJSONRoute(r, prefix)
////		})
////	})
////}
//
//func addHomeRoute(r martini.Router) {
//	r.Get("/", GoHome)
//}
//
//func addUserRoutes(r martini.Router) {
//	r.Get("/user/:uid", GoUser)
//}
//
////func addNewActionRoute(r martini.Router, prefix string) {
////	r.Get("/new",
////		user.RequireUser,
////		gtype.SetTypes,
////		NewAction,
////	).Name(prefix + "_new")
////}
//
////func addCreateRoute(r martini.Router, prefix string) {
////	r.Post("",
////		user.RequireUser,
////		gtype.SetTypes,
////		Create,
////	).Name(prefix + "_create")
////}
////
////func addShowRoute(r martini.Router, prefix string) {
////	r.Get("/:gid",
////		Fetch,
////		mlog.Get,
////		Show,
////	).Name(prefix + "_show")
////}
//
////func addJSONRoute(r martini.Router, prefix string) {
////	r.Get("/:gid/json",
////		Fetch,
////		mlog.Get,
////		JSON,
////	).Name(prefix + "_json")
////}
////
////func addUpdateRoute(r martini.Router, prefix string) {
////	r.Put("/:gid",
////		user.RequireUser,
////		Fetch,
////		CurrentPlayerOrAdmin,
////		Update,
////	).Name(prefix + "_update")
////}
////
////func addAcceptRoute(r martini.Router, prefix string) {
////	r.Put("/:gid/accept",
////		user.RequireUser,
////		Fetch,
////		Accept,
////	).Name(prefix + "_accept")
////}
////
////func addDropRoute(r martini.Router, prefix string) {
////	r.Put("/:gid/drop",
////		user.RequireLogin,
////		Fetch,
////		Drop,
////	).Name(prefix + "_drop")
////}
//
////func addIndexRoute(r martini.Router, prefix string) {
////	r.Get("/:status",
////		user.RequireUser,
////		gtype.SetTypes,
////		SetStatus,
////		SetType,
////		Index,
////	).Name(prefix + "_index")
////
////	r.Post("/:status/json",
////		user.RequireUser,
////		gtype.SetTypes,
////		SetType,
////		GetFiltered,
////		JSONIndexAction,
////	).Name(prefix + "_index_json")
////
////	r.Get("",
////		user.RequireUser,
////		gtype.SetTypes,
////		SetType,
////		Index,
////	).Name(prefix + "_index")
////
////	r.Post("/json",
////		user.RequireUser,
////		gtype.SetTypes,
////		SetType,
////		GetFiltered,
////		JSONIndexAction,
////	).Name(prefix + "_index_json")
////}
//
//func addAddMessageRoute(r martini.Router, prefix string) {
//	r.Put("/:gid/addmessage",
//		user.RequireLogin,
//		Fetch,
//		mlog.Get,
//		mlog.AddMessage,
//	).Name(prefix + "_addmessage")
//}
//
//func addAdminShowRoute(r martini.Router, prefix string) {
//	r.Get("/:gid",
//		user.RequireLogin,
//		user.RequireAdmin,
//		Fetch,
//		mlog.Get,
//		setAdmin,
//		Show,
//	).Name(prefix + "_show")
//}
//
//func addAdminUpdateRoute(r martini.Router, prefix string) {
//	r.Put("/:gid",
//		user.RequireLogin,
//		user.RequireAdmin,
//		Fetch,
//		setAdmin,
//		Update,
//	).Name(prefix + "_update")
//}
//
//func addAdminJSONRoute(r martini.Router, prefix string) {
//	r.Get("/:gid/json",
//		user.RequireLogin,
//		user.RequireAdmin,
//		Fetch,
//		mlog.Get,
//		setAdmin,
//		JSON,
//	).Name(prefix + "_json")
//}
