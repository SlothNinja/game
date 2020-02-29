package game

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/SlothNinja/log"
	"github.com/SlothNinja/restful"
	"github.com/SlothNinja/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
)

func Index(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Debugf("Entering")
		defer log.Debugf("Exiting")

		gs := GamersFrom(c)
		switch status := StatusFrom(c); status {
		case Recruiting:
			c.HTML(http.StatusOK, "shared/invitation_index", gin.H{
				"Context":   c,
				"VersionID": appengine.VersionID(c),
				"CUser":     user.CurrentFrom(c),
				"Games":     gs,
			})
		default:
			c.HTML(http.StatusOK, "shared/multi_games_index", gin.H{
				"Context":   c,
				"VersionID": appengine.VersionID(c),
				"CUser":     user.CurrentFrom(c),
				"Games":     gs,
				"Status":    status,
			})
		}
	}
}

//func Index(c *gin.Context) {
//	switch {
//	case ctx.Data["Status"] == Recruiting:
//		render.HTML(http.StatusOK, "shared/invitation_index", ctx.Data)
//	case ctx.Data["Type"] == gType.All:
//		render.HTML(http.StatusOK, "shared/multi_games_index", ctx.Data)
//	default:
//		render.HTML(http.StatusOK, "shared/games_index", ctx.Data)
//	}
//}

//func GoHome(ctx *restful.Context, render render.Render, routes martini.Routes) {
//	render.Redirect("http://localhost:8080", http.StatusSeeOther)
//}
//
//func GoUser(ctx *restful.Context, render render.Render, routes martini.Routes, params martini.Params) {
//	render.Redirect("http://localhost:8080/user/"+params["uid"], http.StatusSeeOther)
//}

//func NewAction(ctx *restful.Context, render render.Render, routes martini.Routes, params martini.Params) {
//	t := gType.ToType[ctx.Prefix]
//	g := factories[t](ctx)
//	ctx.Data["Game"] = g
//	route := routes.URLFor(ctx.Prefix+"_games_index", "recruiting")
//	if err := g.FromParams(ctx, t); err != nil {
//		ctx.Errorf("SN/Game/Controller#NewAction NewAction Error: %s", err)
//		render.Redirect(route, http.StatusSeeOther)
//		return
//	}
//	render.HTML(http.StatusOK, ctx.Prefix+"/new", ctx.Data)
//}

//func Create(ctx *restful.Context, render render.Render, routes martini.Routes, params martini.Params) {
//	t := gType.ToType[ctx.Prefix]
//	g := factories[t](ctx)
//	ctx.Data["Game"] = g
//	route := routes.URLFor(ctx.Prefix+"_games_index", "recruiting")
//	if err := g.FromParams(ctx, t); err != nil {
//		ctx.Errorf("SN/Game/Controller#Create g.FromParams Error: %s", err)
//		render.Redirect(route, http.StatusSeeOther)
//		return
//	}
//
//	m := mlog.New(ctx)
//
//	if err := datastore.RunInTransaction(ctx, func(tc context.Context) error {
//		if err := datastore.Put(tc, g); err != nil {
//			return err
//		}
//
//		m.ID = g.GetHeader().ID
//
//		return datastore.Put(tc, m)
//	}, nil); err != nil {
//		ctx.Errorf("SN/Game/Controller#Create datastore.RunInTransaction Error: %s", err)
//		render.Redirect(route, http.StatusSeeOther)
//		return
//	}
//
//	ctx.AddNoticef("<div>%s created.</div>", g.GetTitle())
//	render.Redirect(route, http.StatusSeeOther)
//}

//func Accept(ctx *restful.Context, render render.Render, routes martini.Routes) {
//	g := gamer(ctx)
//	if g == nil {
//		ctx.Errorf("Controller#Accept Game Not Found")
//		render.Redirect(routes.URLFor(ctx.Prefix+"_games_index", "recruiting"), http.StatusSeeOther)
//		return
//	}
//
//	h := g.GetHeader()
//	u := user.Current(ctx)
//	ctx.Debugf("prefix: %s", ctx.Prefix)
//	route := routes.URLFor(ctx.Prefix+"_games_index", "recruiting")
//	if ctx.Req.Form.Get("action") == "get-dialog" {
//		render.HTML(http.StatusOK, ctx.Prefix+"/accept_dialog", ctx.Data)
//		return
//	}
//
//	start, err := g.Accept(ctx, u)
//	if err != nil {
//		ctx.Errorf(err.Error())
//		ctx.AddErrorf(err.Error())
//		render.Redirect(route, http.StatusSeeOther)
//		return
//	}
//
//	if err := datastore.RunInTransaction(ctx, func(tc context.Context) error {
//		if start {
//			if err := g.Start(); err != nil {
//				return err
//			}
//		}
//
//		return datastore.Put(tc, g)
//	}, nil); err != nil {
//		ctx.Errorf("Controller#Accept Error: %s", err)
//		ctx.AddErrorf("Controller#Accept Error: %s", err)
//	} else {
//		ctx.AddNoticef("%s accepted invitation for Game #%d: %s.", u.Name, h.ID, h.Title)
//		if start {
//			g.SendTurnNotificationsToCurrentPlayers()
//		}
//	}
//	render.Redirect(route, http.StatusSeeOther)
//}
//
//func Drop(ctx *restful.Context, render render.Render, routes martini.Routes) {
//	g := gamer(ctx)
//	if g == nil {
//		ctx.Errorf("Controller#Drop Game Not Found")
//		render.Redirect(routes.URLFor(ctx.Prefix+"_games_index", "recruiting"), http.StatusSeeOther)
//		return
//	}
//
//	h := g.GetHeader()
//	u := user.Current(ctx)
//	route := routes.URLFor(ctx.Prefix+"_games_index", "recruiting")
//	if err := g.Drop(u); err != nil {
//		ctx.Errorf(err.Error())
//		ctx.AddErrorf(err.Error())
//		render.Redirect(route, http.StatusSeeOther)
//		return
//	}
//
//	if err := datastore.RunInTransaction(ctx, func(tc context.Context) error {
//		return datastore.Put(tc, g)
//	}, nil); err != nil {
//		ctx.Errorf("Controller#Drop Error: %s", err)
//		ctx.AddErrorf("Controller#Drop Error: %s", err)
//	} else {
//		ctx.AddNoticef("%s dropped from Game #%d: %s.", u.Name, h.ID, h.Title)
//	}
//	render.Redirect(route, http.StatusSeeOther)
//}

type ActionType int

const (
	None ActionType = iota
	Save
	SaveAndStatUpdate
	Cache
	UndoAdd
	UndoReplace
	UndoPop
	Undo
	Redo
	Reset
)

const (
	gamerKey  = "Game"
	gamersKey = "Games"
	homePath  = "/"
	adminKey  = "Admin"
)

type Action struct {
	Call func(Gamer) (string, error)
	Type ActionType
}

//func Update(ctx *restful.Context, render render.Render, routes martini.Routes, params martini.Params, form url.Values) {
//	g := gamer(ctx)
//	if g == nil {
//		ctx.Errorf("Controller#Update Game Not Found")
//		render.Redirect(routes.URLFor("home"), http.StatusSeeOther)
//		return
//	}
//	template, actionType, err := g.(hasUpdate).Update(ctx, form)
//	switch {
//	case err != nil && sn.IsVError(err):
//		ctx.AddErrorf("%v", err)
//		ctx.Data["JSON"] = g
//	case err != nil:
//		ctx.Errorf(err.Error())
//		render.Redirect(routes.URLFor("home"), http.StatusSeeOther)
//		return
//	case actionType == Cache:
//		mkey := g.GetHeader().undoKey(ctx)
//		item := memcache.NewItem(ctx, mkey).SetExpiration(time.Minute * 30)
//		v, err := codec.Encode(g)
//		if err != nil {
//			ctx.Errorf("Controller#Update Cache Error: %s", err)
//			render.Redirect(routes.URLFor(ctx.Prefix+"_game_show", params["gid"]), http.StatusSeeOther)
//			return
//		}
//		item.SetValue(v)
//		if err := memcache.Set(ctx, item); err != nil {
//			ctx.Errorf("Controller#Update Cache Error: %s", err)
//			render.Redirect(routes.URLFor(ctx.Prefix+"_game_show", params["gid"]), http.StatusSeeOther)
//			return
//		}
//	case actionType == UndoPop:
//		mkey := g.GetHeader().undoKey(ctx)
//		s := stack(ctx)
//		s.Pop()
//
//		item := memcache.NewItem(ctx, mkey).SetExpiration(time.Minute * 30)
//		v, err := codec.Encode(s)
//		if err != nil {
//			ctx.Errorf("Controller#Update Cache Error: %s", err)
//			render.Redirect(routes.URLFor(ctx.Prefix+"_game_show", params["gid"]), http.StatusSeeOther)
//			return
//		}
//		item.SetValue(v)
//		if err := memcache.Set(ctx, item); err != nil {
//			ctx.Errorf("Controller#Update Cache Error: %s", err)
//			render.Redirect(routes.URLFor(ctx.Prefix+"_game_show", params["gid"]), http.StatusSeeOther)
//			return
//		}
//	case actionType == UndoReplace:
//		mkey := g.GetHeader().undoKey(ctx)
//		entry := undo.NewEntry(form)
//		s := stack(ctx)
//		s.Replace(entry)
//
//		item := memcache.NewItem(ctx, mkey).SetExpiration(time.Minute * 30)
//		v, err := codec.Encode(s)
//		if err != nil {
//			ctx.Errorf("Controller#Update Cache Error: %s", err)
//			render.Redirect(routes.URLFor(ctx.Prefix+"_game_show", params["gid"]), http.StatusSeeOther)
//		}
//		item.SetValue(v)
//		if err := memcache.Set(ctx, item); err != nil {
//			ctx.Errorf("Controller#Update Cache Error: %s", err)
//			render.Redirect(routes.URLFor(ctx.Prefix+"_game_show", params["gid"]), http.StatusSeeOther)
//		}
//	case actionType == UndoAdd:
//		mkey := g.GetHeader().undoKey(ctx)
//		entry := undo.NewEntry(form)
//		s := stack(ctx)
//		s.Push(entry)
//
//		item := memcache.NewItem(ctx, mkey).SetExpiration(time.Minute * 30)
//		v, err := codec.Encode(s)
//		if err != nil {
//			ctx.Errorf("Controller#Update Cache Error: %s", err)
//			render.Redirect(routes.URLFor(ctx.Prefix+"_game_show", params["gid"]), http.StatusSeeOther)
//			return
//		}
//		item.SetValue(v)
//		if err := memcache.Set(ctx, item); err != nil {
//			ctx.Errorf("Controller#Update Cache Error: %s", err)
//			render.Redirect(routes.URLFor(ctx.Prefix+"_game_show", params["gid"]), http.StatusSeeOther)
//			return
//		}
//	case actionType == SaveAndStatUpdate:
//		if err := saveAndUpdateStats(ctx, g); err != nil {
//			ctx.Errorf("%s", err)
//			ctx.AddErrorf("Controller#Update SaveAndStatUpdate Error: %s", err)
//			render.Redirect(routes.URLFor(ctx.Prefix+"_game_show", params["gid"]), http.StatusSeeOther)
//			return
//		}
//	case actionType == Save:
//		if err := save(ctx, g); err != nil {
//			ctx.Errorf("%s", err)
//			ctx.AddErrorf("Controller#Update Save Error: %s", err)
//			render.Redirect(routes.URLFor(ctx.Prefix+"_game_show", params["gid"]), http.StatusSeeOther)
//			return
//		}
//	case actionType == Reset:
//		mkey := g.GetHeader().undoKey(ctx)
//		if err := memcache.Delete(ctx, mkey); err != nil && err != memcache.ErrCacheMiss {
//			ctx.Errorf("Controller#Reset Error: %s", err)
//			render.Redirect(routes.URLFor(ctx.Prefix+"_game_show", params["gid"]), http.StatusSeeOther)
//			return
//		}
//		if sdd, ok := ctx.Data["Game"].(SetDefaultDialoger); ok {
//			sdd.SetDefaultDialog(ctx)
//		}
//		if sdu, ok := ctx.Data["Game"].(SetDefaultUIer); ok {
//			sdu.SetDefaultUI(ctx)
//		}
//		// Fetch(ctx, routes, params)
//		//		switch err := gaelic.Get(ctx, g.GetKey(), g); {
//		//		case err != nil:
//		//			ctx.AddErrorf(err.Error())
//		//			render.Redirect(routes.URLFor("home"), http.StatusSeeOther)
//		//			return
//		//		case g == nil:
//		//			ctx.AddErrorf("Unable to get game for id: %v", params["gid"])
//		//			render.Redirect(routes.URLFor("home"), http.StatusSeeOther)
//		//			return
//		//		}
//		//
//		//		ctx.Data["Game"] = g
//	case actionType == Redo:
//		c := stack(ctx).Forward()
//		h := g.GetHeader()
//		mkey := h.undoKey(ctx)
//		item := memcache.NewItem(ctx, mkey).SetExpiration(time.Minute * 30)
//		v, err := codec.Encode(c)
//		if err != nil {
//			ctx.Errorf("Controller#Update Undo Error: %s", err)
//			render.Redirect(routes.URLFor(ctx.Prefix+"_game_show", params["gid"]), http.StatusSeeOther)
//		}
//		item.SetValue(v)
//		if err := memcache.Set(ctx, item); err != nil {
//			ctx.Errorf("Controller#Update Undo Error: %s", err)
//			render.Redirect(routes.URLFor(ctx.Prefix+"_game_show", params["gid"]), http.StatusSeeOther)
//		}
//	case actionType == Undo && g.GetType().MultiUndo():
//		mkey := g.GetHeader().undoKey(ctx)
//		c := stack(ctx).Back()
//		item := memcache.NewItem(ctx, mkey).SetExpiration(time.Minute * 30)
//		v, err := codec.Encode(c)
//		if err != nil {
//			ctx.Errorf("Controller#Update Undo Error: %s", err)
//			render.Redirect(routes.URLFor(ctx.Prefix+"_game_show", params["gid"]), http.StatusSeeOther)
//		}
//		item.SetValue(v)
//		if err := memcache.Set(ctx, item); err != nil {
//			ctx.Errorf("Controller#Update Undo Error: %s", err)
//			render.Redirect(routes.URLFor(ctx.Prefix+"_game_show", params["gid"]), http.StatusSeeOther)
//		}
//		if sdd, ok := ctx.Data["Game"].(SetDefaultDialoger); ok {
//			sdd.SetDefaultDialog(ctx)
//		}
//		if sdu, ok := ctx.Data["Game"].(SetDefaultUIer); ok {
//			sdu.SetDefaultUI(ctx)
//		}
//	case actionType == Undo && !g.GetType().MultiUndo():
//		mkey := g.GetHeader().undoKey(ctx)
//		if err := memcache.Delete(ctx, mkey); err != nil && err != memcache.ErrCacheMiss {
//			ctx.Errorf("Controller#Undo Error: %s", err)
//			render.Redirect(routes.URLFor(ctx.Prefix+"_game_show", params["gid"]), http.StatusSeeOther)
//		}
//	}
//
//	switch jData, ok := ctx.Data["JSON"]; {
//	case ok && template == "json":
//		render.JSON(http.StatusOK, jData)
//	case template == "":
//		render.Redirect(routes.URLFor(ctx.Prefix+"_game_show", params["gid"]), http.StatusSeeOther)
//	default:
//		render.HTML(http.StatusOK, template, ctx.Data)
//	}
//}
//
//func save(ctx *restful.Context, g Gamer) error {
//	return datastore.RunInTransaction(ctx, func(tc context.Context) error {
//		oldG := factories[g.GetType()](ctx)
//		if ok := datastore.PopulateKey(oldG, datastore.KeyForObj(tc, g)); !ok {
//			return fmt.Errorf("Unable to populate game with key.")
//		}
//
//		if err := datastore.Get(tc, oldG); err != nil {
//			return err
//		}
//
//		if oldG.GetHeader().UpdatedAt != g.GetHeader().UpdatedAt {
//			return fmt.Errorf("Game state changed unexpectantly.  Try again.")
//		}
//
//		if err := datastore.Put(tc, g); err != nil {
//			return err
//		}
//		return memcache.Delete(ctx, g.GetHeader().undoKey(ctx))
//	}, nil)
//}
//
//func saveAndUpdateStats(ctx *restful.Context, g Gamer) error {
//	return datastore.RunInTransaction(ctx, func(tc context.Context) error {
//		oldG := factories[g.GetType()](ctx)
//		if ok := datastore.PopulateKey(oldG.GetHeader(), datastore.KeyForObj(tc, g.GetHeader())); !ok {
//			return fmt.Errorf("Unable to populate game with key.")
//		}
//		if err := datastore.Get(tc, oldG.GetHeader()); err != nil {
//			return err
//		}
//
//		ctx.Debugf("Old: %v New: %v", oldG.GetHeader().UpdatedAt, g.GetHeader().UpdatedAt)
//		if oldG.GetHeader().UpdatedAt != g.GetHeader().UpdatedAt {
//			return fmt.Errorf("Game state changed unexpectantly.  Try again.")
//		}
//
//		s := stats.Fetched(ctx)
//		if s == nil {
//			return fmt.Errorf("Unable to find player stats. Try again.")
//		}
//
//		entities := []interface{}{s, g}
//		if err := datastore.Put(tc, entities); err != nil {
//			return err
//		}
//		return memcache.Delete(tc, g.GetHeader().undoKey(ctx))
//	}, nil)
//}

//type SetDefaultDialoger interface {
//	SetDefaultDialog(*restful.Context) error
//}
//
//type SetDefaultUIer interface {
//	SetDefaultUI(*restful.Context) error
//}

//func Show(ctx context.Context, render render.Render) {
//	c := restful.GinFrom(ctx)
//	c.HTML(http.StatusOK, showPath(ctx), gin.H{})
//}

//func JSON(ctx *restful.Context, render render.Render) {
//	if sdd, ok := ctx.Data["Game"].(SetDefaultDialoger); ok {
//		sdd.SetDefaultDialog(ctx)
//	}
//	if sdu, ok := ctx.Data["Game"].(SetDefaultUIer); ok {
//		sdu.SetDefaultUI(ctx)
//	}
//	render.JSON(http.StatusOK, ctx.Data["Game"])
//}

func NotCurrentPlayerOrAdmin(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch g := GamerFrom(c); {
		case g == nil:
			restful.AddErrorf(c, "Did not load game from database.")
			c.Redirect(http.StatusSeeOther, homePath)
		case g.GetHeader().CUserIsCPlayerOrAdmin(c):
			c.Redirect(http.StatusSeeOther, showPath(c, prefix, c.Param("gid")))
		}
	}
}

func CurrentPlayerOrAdmin(prefix, idParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch g := GamerFrom(c); {
		case g == nil:
			restful.AddErrorf(c, "Did not load game from database.")
			c.Redirect(http.StatusSeeOther, homePath)
		case !g.GetHeader().CUserIsCPlayerOrAdmin(c):
			c.Redirect(http.StatusSeeOther, showPath(c, prefix, c.Param(idParam)))
		}
	}
}

func showPath(c *gin.Context, prefix, id string) string {
	return fmt.Sprintf("/%s/game/%s/show", prefix, id)
}

func GamerFrom(c *gin.Context) (g Gamer) {
	g, _ = c.Value(gamerKey).(Gamer)
	return
}

func WithGamer(c *gin.Context, g Gamer) *gin.Context {
	c.Set(gamerKey, g)
	return c
}

func GamersFrom(c *gin.Context) (gs Gamers) {
	gs, _ = c.Value(gamersKey).(Gamers)
	return
}

func withGamers(c *gin.Context, gs Gamers) *gin.Context {
	c.Set(gamersKey, gs)
	return c
}

//func stack(c *restful.Context) *undo.Stack {
//	c1, ok1 := ctx.Data["Undo"]
//	if c1 == nil || !ok1 {
//		return nil
//	}
//
//	c2, ok2 := c1.(*undo.Stack)
//	if !ok2 {
//		return nil
//	}
//	return c2
//}
//
//func IsLastAction(ctx *restful.Context, a string) bool {
//	s := stack(ctx)
//	return s.Top() != nil && s.Top().Get("action") == a
//}

type dbState interface {
	DBState()
}

//func Fetch(ctx *restful.Context, render render.Render, routes martini.Routes, params martini.Params, form url.Values) {
//	// create Gamer
//	g, t, err := create(ctx, params["gid"])
//	if err != nil {
//		render.Redirect(routes.URLFor("home"), http.StatusSeeOther)
//	}
//
//	switch action := form.Get("action"); {
//	case action == "reset":
//		// pull from memcache/datastore
//		// same as undo & !MultiUndo
//		fallthrough
//	case action == "undo" && !t.MultiUndo():
//		// pull from memcache/datastore
//		if err := dsGet(ctx, g); err != nil {
//			render.Redirect(routes.URLFor("home"), http.StatusSeeOther)
//			return
//		}
//	case action == "undo" && t.MultiUndo():
//		// pull from datastore then playback command stack - 1
//		if err := dsGet(ctx, g); err != nil {
//			render.Redirect(routes.URLFor("home"), http.StatusSeeOther)
//			return
//		}
//		if err := playBack(ctx, g, -1); err != nil {
//			render.Redirect(routes.URLFor("home"), http.StatusSeeOther)
//			return
//		}
//	case action == "redo" && t.MultiUndo():
//		// pull from datastore and playback command stack + 1
//		if err := dsGet(ctx, g); err != nil {
//			render.Redirect(routes.URLFor("home"), http.StatusSeeOther)
//			return
//		}
//		if err := playBack(ctx, g, 1); err != nil {
//			render.Redirect(routes.URLFor("home"), http.StatusSeeOther)
//			return
//		}
//	case !t.MultiUndo() && user.Current(ctx) != nil:
//		// pull from memcache and return if successful; otherwise pull from datastore
//		if err := mcGet(ctx, g); err == nil {
//			return
//		}
//		if err := dsGet(ctx, g); err != nil {
//			render.Redirect(routes.URLFor("home"), http.StatusSeeOther)
//			return
//		}
//	default:
//		// Should be MultiUndo.
//		// Pull from memcache/datastore and playback command stack.
//		if err := dsGet(ctx, g); err != nil {
//			render.Redirect(routes.URLFor("home"), http.StatusSeeOther)
//			return
//
//		}
//		if user.Current(ctx) != nil {
//			if err := playBack(ctx, g, 0); err != nil {
//				render.Redirect(routes.URLFor("home"), http.StatusSeeOther)
//				return
//			}
//		}
//	}
//}

//// create appropriate gamer from factory based on type in url prefix.
//func create(ctx *restful.Context, gid string) (Gamer, gType.Type, error) {
//	t := gType.ToType[ctx.Prefix]
//	id, err := strconv.ParseInt(gid, 10, 64)
//	if err != nil {
//		ctx.AddErrorf(err.Error())
//		return nil, gType.NoType, err
//	}
//	g := factories[t](ctx)
//	h := g.GetHeader()
//	h.ID = id
//
//	// provide debug hook
//	if stater, ok := g.(dbState); ok {
//		stater.DBState()
//	}
//	return g, t, nil
//}
//
//// pull temporary game state from memcache.  Note may be different from value stored in datastore.
//func mcGet(ctx *restful.Context, g Gamer) error {
//	mkey := g.GetHeader().undoKey(ctx)
//	item, err := memcache.GetKey(ctx, mkey)
//	if err != nil {
//		return err
//	}
//
//	if err := codec.Decode(g, item.Value()); err != nil {
//		return err
//	}
//	if afterCacher, ok := g.(restful.AfterCacher); ok {
//		afterCacher.AfterCache()
//	}
//	ctx.Data["Game"] = g
//	ctx.Data["ColorMap"] = g.ColorMapFor(user.Current(ctx))
//	return nil
//}
//
//// pull game state from memcache/datastore.  returned memcache should be same as datastore.
//func dsGet(ctx *restful.Context, g Gamer) error {
//	switch err := datastore.Get(ctx, g); {
//	case err != nil:
//		ctx.AddErrorf(err.Error())
//		return err
//	case g == nil:
//		err := fmt.Errorf("Unable to get game for id: %v", g.GetHeader().ID)
//		ctx.AddErrorf(err.Error())
//		return err
//	}
//
//	ctx.Data["Game"] = g
//	ctx.Data["ColorMap"] = g.ColorMapFor(user.Current(ctx))
//	return nil
//}
//
//// playback command stack up to current level but adjusted by adj
//func playBack(ctx *restful.Context, g Gamer, adj int) error {
//	stack := new(undo.Stack)
//	ctx.Data["Undo"] = stack
//	mkey := g.GetHeader().undoKey(ctx)
//	item, err := memcache.GetKey(ctx, mkey)
//	if err != nil {
//		return err
//	}
//	if err := codec.Decode(stack, item.Value()); err == nil {
//		stop := stack.Current + adj
//		switch {
//		case stop < 0:
//			stop = 0
//		case stop > stack.Count():
//			stop = stack.Count()
//		}
//		for i := 0; i < stop; i++ {
//			entry := stack.Entries[i]
//			if _, _, err := g.(hasUpdate).Update(ctx, entry.Values); err != nil {
//				ctx.AddErrorf("Unexpected error.  Reset turn and try again.")
//				ctx.Errorf("Fetch Error: %#v", err)
//				return err
//			}
//		}
//	}
//	return nil
//}

//func setAdmin(ctx *restful.Context) {
//	ctx.Data["Admin"] = true
//}

func AdminFrom(c *gin.Context) (b bool) {
	b, _ = c.Value(adminKey).(bool)
	return
}

func WithAdmin(c *gin.Context, b bool) {
	c.Set(adminKey, b)
}

func SetAdmin(b bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		WithAdmin(c, b)
	}
}

func (h *Header) undoKey(c *gin.Context) (k string) {
	if cu := user.CurrentFrom(c); cu != nil {
		k = fmt.Sprintf("%s/uid-%d", h.Key, cu.ID)
		// k = fmt.Sprintf("%s/uid-%d", dscache.MakeMemcacheKey(0, datastore.KeyForObj(ctx, h)), cu.ID)
	}
	return
}

func (h *Header) UndoKey(c *gin.Context) string {
	return h.undoKey(c)
}

type jGamesIndex struct {
	Data            []*jHeader `json:"data"`
	Draw            int        `json:"draw"`
	RecordsTotal    int64      `json:"recordsTotal"`
	RecordsFiltered int64      `json:"recordsFiltered"`
}

type jHeader struct {
	ID          int64         `json:"id"`
	Type        template.HTML `json:"type"`
	Title       template.HTML `json:"title"`
	Creator     template.HTML `json:"creator"`
	Players     template.HTML `json:"players"`
	NumPlayers  template.HTML `json:"numPlayers"`
	OptString   template.HTML `json:"optString"`
	Progress    template.HTML `json:"progress"`
	Round       int           `json:"round"`
	UpdatedAt   time.Time     `json:"updatedAt"`
	LastUpdated template.HTML `json:"lastUpdated"`
	Public      template.HTML `json:"public"`
	Actions     template.HTML `json:"actions"`
	Status      Status        `json:"status"`
}

func JSONIndexAction(c *gin.Context) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	if data, err := toGameTable(c); err != nil {
		c.JSON(http.StatusOK, fmt.Sprintf("%v", err))
	} else {
		c.JSON(http.StatusOK, data)
	}
}

func toGameTable(c *gin.Context) (*jGamesIndex, error) {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	gs := GamersFrom(c)
	table := new(jGamesIndex)
	l := len(gs)
	table.Data = make([]*jHeader, l)
	for i, g := range gs {
		h := g.GetHeader()
		table.Data[i] = &jHeader{
			ID:          g.GetHeader().ID(),
			Type:        template.HTML(h.Type.String()),
			Title:       titleLink(g),
			Creator:     user.LinkFor(h.CreatorID, h.CreatorName),
			Players:     h.PlayerLinks(c),
			NumPlayers:  template.HTML(fmt.Sprintf("%d / %d", h.AcceptedPlayers(), h.NumPlayers)),
			Round:       h.Round,
			OptString:   template.HTML(h.OptString),
			Progress:    template.HTML(h.Progress),
			UpdatedAt:   h.UpdatedAt,
			LastUpdated: template.HTML(restful.LastUpdated(time.Time(h.UpdatedAt))),
			Public:      publicPrivate(g),
			Actions:     actionButtons(c, g),
			Status:      h.Status,
		}
	}

	if draw, err := strconv.Atoi(c.PostForm("draw")); err != nil {
		return nil, err
	} else {
		table.Draw = draw
	}
	table.RecordsTotal = countFrom(c)
	table.RecordsFiltered = countFrom(c)
	return table, nil
}

func publicPrivate(g Gamer) template.HTML {
	h := g.GetHeader()
	if h.Private() {
		return template.HTML("Private")
	} else {
		return template.HTML("Public")
	}
}

func titleLink(g Gamer) template.HTML {
	h := g.GetHeader()
	return template.HTML(fmt.Sprintf(`
		<div><a href="/%s/game/show/%d">%s</a></div>
		<div style="font-size:.7em">%s</div>`, h.Type.IDString(), h.ID, h.Title, h.OptString))
}

func actionButtons(c *gin.Context, g Gamer) template.HTML {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	h := g.GetHeader()
	switch h.Status {
	case Running:
		t := h.Type.IDString()
		if g.GetHeader().CUserIsCPlayerOrAdmin(c) {
			return template.HTML(fmt.Sprintf(`<a class="mybutton" href="/%s/game/show/%d">Play</a>`, t, h.ID))
		} else {
			return template.HTML(fmt.Sprintf(`<a class="mybutton" href="/%s/game/show/%d">Show</a>`, t, h.ID))
		}
	case Recruiting:
		t := h.Type.IDString()
		switch cu := g.CurrentUser(); {
		case g.CanAdd(cu):
			if g.Private() {
				return template.HTML(fmt.Sprintf(`
	<div id="dialog-%d" title="Game %d">
		<form class="top-padding" action="/%s/game/accept/%d" method="post">
			<input id="password" name="password" type="text" value="Please Enter Password" />
			<div>
				&nbsp;
			</div>
			<div class="top-padding center" >
				<input type="submit" value="Accept" class="mybutton" />
			</div>
		</form>
	</div>
	<div id="opener-%d" class="mybutton">Accept</div>
	<script>
		$('#dialog-%d').dialog({autoOpen: false, modal: true});
		$('#opener-%d').click(function() {
			$('#dialog-%d').dialog('open');
		});
	</script>`, h.ID, h.ID, h.Stub(), h.ID, h.ID, h.ID, h.ID, h.ID))
				//				return template.HTML(fmt.Sprintf(`
				//				<form class="myForm" action="/%s/game/accept/%d" method="post">
				//					<input name="_method" type="hidden" value="PUT" />
				//					<input name="action" type="hidden" value="get-dialog" />
				//					<input type="submit" value="Accept" class="mybutton" />
				//				</form>`, t, h.ID))
			} else {
				return template.HTML(fmt.Sprintf(`
				<form method="post" action="/%s/game/accept/%d">
					<input name="_method" type="hidden" value="PUT" />
					<input id="user_id" name="user[id]" type="hidden" value="%v">
					<input id="accept-%d" class="mybutton" type="submit" value="Accept" />
				</form>`, t, h.ID, cu.ID, h.ID))
			}
		case g.CanDropout(cu):
			return template.HTML(fmt.Sprintf(`
				<form method="post" action="/%s/game/drop/%d">
					<input name="_method" type="hidden" value="PUT" />
					<input id="user_id" name="user[id]" type="hidden" value="%v">
					<input id="drop-%d" class="mybutton" type="submit" value="Drop" />
				</form>`, t, h.ID, cu.ID, h.ID))
		default:
			return ""
		}
	default:
		return ""
	}
}
