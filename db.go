package game

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/log"
	gtype "github.com/SlothNinja/type"
	"github.com/gin-gonic/gin"
)

const (
	statusKey = "Status"
	countKey  = "Count"
	NoCount   = -1
)

func getAllQuery(c *gin.Context) *datastore.Query {
	return datastore.NewQuery("Game").Ancestor(GamesRoot(c))
}

func getFiltered(c *gin.Context, status, sid, start, length string, t gtype.Type) (Gamers, int64, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	dsClient, err := datastore.NewClient(c, "")
	if err != nil {
		return nil, 0, err
	}

	q := getAllQuery(c).
		KeysOnly()

	if status != "" {
		st := ToStatus[strings.ToLower(status)]
		q = q.Filter("Status=", int(st))
		WithStatus(c, st)
	}

	if sid != "" {
		if id, err := strconv.Atoi(sid); err == nil {
			q = q.Filter("UserIDS=", id)
		}
	}

	if t != gtype.All {
		q = q.Filter("Type=", int(t)).
			Order("-UpdatedAt")
	} else {
		q = q.Order("-UpdatedAt")
	}

	cnt, err := dsClient.Count(c, q)
	if err != nil {
		log.Errorf("sn/game#GetFiltered q.Count Error: %s", err)
		return nil, 0, err
	}

	if start != "" {
		if st, err := strconv.ParseInt(start, 10, 32); err == nil {
			q = q.Offset(int(st))
		}
	}

	if length != "" {
		if l, err := strconv.ParseInt(length, 10, 32); err == nil {
			q = q.Limit(int(l))
		}
	}

	ks, err := dsClient.GetAll(c, q, nil)
	if err != nil {
		log.Errorf("getFiltered GetAll Error: %s", err)
		return nil, 0, err
	}

	l := len(ks)
	gs := make([]Gamer, l)
	hs := make([]*Header, l)
	for i := range gs {
		var ok bool
		if t == gtype.All {
			k := strings.ToLower(ks[i].Parent.Kind)
			if t, ok = gtype.ToType[k]; !ok {
				err = fmt.Errorf("Unknown Game Type For: %s", k)
				log.Errorf(err.Error())
				return nil, 0, err
			}
		}
		gs[i] = factories[t](c)
		hs[i] = gs[i].GetHeader()
		// if ok := datastore.PopulateKey(hs[i], ks[i]); !ok {
		// 	err = fmt.Errorf("Unable to populate header with key.")
		// 	log.Errorf(err.Error())
		// 	return
		// }
	}

	err = dsClient.GetMulti(c, ks, hs)
	if err != nil {
		log.Errorf("SN/Game#GetFiltered datastore.Get Error: %s", err)
		return nil, 0, err
	}

	for i := range hs {
		if err = hs[i].AfterLoad(gs[i]); err != nil {
			log.Errorf("SN/Game#GetFiltered h.AfterLoad Error: %s", err)
			return nil, 0, err
		}
	}

	return gs, int64(cnt), nil
}

//func SetStatus(c *gin.Context) {
//	ctx := restful.ContextFrom(c)
//	log.Debugf(ctx, "Entering")
//	defer log.Debugf(ctx, "Exiting")
//
//	stat := c.Param("status")
//	status := ToStatus[stat]
//	WithStatus(c, status)
//}

func WithStatus(c *gin.Context, s Status) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	c.Set(statusKey, s)
}

func StatusFrom(c *gin.Context) (s Status) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	if s = ToStatus[strings.ToLower(c.Param("status"))]; s != NoStatus {
		WithStatus(c, s)
	} else {
		s, _ = c.Value(statusKey).(Status)
	}
	return
}

func withCount(c *gin.Context, cnt int64) *gin.Context {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	c.Set(countKey, cnt)
	return c
}

func countFrom(c *gin.Context) (cnt int64) {
	cnt, _ = c.Value(countKey).(int64)
	return
}

//func SetType(ctx *gin.Context) {
//	ctx = gtype.WithType(ctx, gtype.All)
//	prefix := restful.PrefixFrom(ctx)
//	if t, ok := gtype.ToType[prefix]; ok {
//		ctx = gtype.WithType(ctx, t)
//		return
//	}
//}

//func GetAll(c *gin.Context) {
//	q := getAllQuery(c).Order("-UpdatedAt").KeysOnly(true)
//
//	if stat := c.Param("status"); stat != "" {
//		status := ToStatus[stat]
//		q = q.Eq("Status", status)
//		c = WithStatus(c, status)
//	}
//
//	if sid := c.Param("uid"); sid != "" {
//		if id, err := strconv.Atoi(sid); err == nil {
//			q = q.Eq("UserIDS", id)
//		}
//	}
//
//	var limit int32 = 100
//	prefix := restful.PrefixFrom(c)
//	if t, ok := gtype.ToType[prefix]; ok {
//		c = gtype.WithType(c, t)
//		q = q.Eq("Type", t).Order("-UpdatedAt").Limit(limit)
//	} else {
//		q = q.Order("-UpdatedAt").Limit(limit)
//	}
//
//	var ks []*datastore.Key
//	ctx := restful.ContextFrom(c)
//	if err := datastore.GetAll(ctx, q, &ks); err != nil {
//		log.Errorf(c, "SN/Game#All Error: %s", err)
//		c.Redirect(http.StatusSeeOther, homePath)
//		return
//	}
//
//	length := len(ks)
//	gs := make([]Gamer, length)
//	for i := range gs {
//		var (
//			t  gtype.Type
//			ok bool
//		)
//		if t, ok = gtype.ToType[prefix]; !ok {
//			k := strings.ToLower(ks[i].Parent().Kind())
//			if t, ok = gtype.ToType[k]; !ok {
//				log.Errorf(c, "Unknown Game Type For: %s", k)
//				c.Redirect(http.StatusSeeOther, homePath)
//				return
//			}
//		}
//		gs[i] = factories[t](c)
//		if ok := datastore.PopulateKey(gs[i], ks[i]); !ok {
//			log.Errorf(c, "Unable to populate gamer with key")
//			c.Redirect(http.StatusSeeOther, homePath)
//			return
//		}
//	}
//
//	if err := datastore.Get(ctx, gs); err != nil {
//		log.Errorf(c, "SN/Game#All gaelic.GetMulti Error: %s", err)
//		c.Redirect(http.StatusSeeOther, homePath)
//		return
//	}
//
//	WithGamers(c, gs)
//}

func GetFiltered(t gtype.Type) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		gs, cnt, err := getFiltered(c, c.Param("status"), c.Param("uid"), c.PostForm("start"), c.PostForm("length"), t)

		if err != nil {
			log.Errorf(err.Error())
			c.Redirect(http.StatusSeeOther, homePath)
			c.Abort()
		}
		withGamers(withCount(c, cnt), gs)
	}
}

func GetRunning(c *gin.Context) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	gs, cnt, err := getFiltered(c, c.Param("status"), "", "", "", gtype.All)

	if err != nil {
		log.Errorf(err.Error())
		c.Redirect(http.StatusSeeOther, homePath)
		c.Abort()
	}
	withGamers(withCount(c, cnt), gs)
}

//	q := datastore.NewQuery("Game").Ancestor(GamesRoot(c)).Eq("Status", Running).KeysOnly(true)
//
//	var ks []*datastore.Key
//	ctx := restful.ContextFrom(c)
//	if err := datastore.GetAll(ctx, q, &ks); err != nil {
//		log.Errorf(c, "SN/Game#All Error: %s", err)
//		c.Redirect(http.StatusSeeOther, homePath)
//		return
//	}
//
//	length := len(ks)
//	gs := make([]Gamer, length)
//	for i := range gs {
//		var (
//			t  gtype.Type
//			ok bool
//		)
//		k := strings.ToLower(ks[i].Parent().Kind())
//		if t, ok = gtype.ToType[k]; !ok {
//			log.Errorf(c, "Unknown Game Type For: %s", k)
//			c.Redirect(http.StatusSeeOther, homePath)
//			return
//		}
//		gs[i] = factories[t](c)
//		if ok := datastore.PopulateKey(gs[i], ks[i]); !ok {
//			log.Errorf(c, "Unable to populate gamer with key.")
//			c.Redirect(http.StatusSeeOther, homePath)
//			return
//		}
//	}
//
//	if err := datastore.Get(ctx, gs); err != nil {
//		log.Errorf(c, "SN/Game#All datastore.Get Error: %s", err)
//		c.Redirect(http.StatusSeeOther, homePath)
//		return
//	}
//
//	WithGamers(c, gs)
//}
