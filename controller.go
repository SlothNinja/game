package game

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/log"
	"github.com/SlothNinja/restful"
	"github.com/SlothNinja/sn"
	gType "github.com/SlothNinja/type"
	"github.com/SlothNinja/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

func (client *Client) Index(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		client.Log.Debugf("Entering")
		defer client.Log.Debugf("Exiting")

		gs := GamersFrom(c)
		cu, err := client.User.Current(c)
		if err != nil {
			client.Log.Debugf(err.Error())
		}

		status := StatusFrom(c)
		switch status {
		case Recruiting:
			c.HTML(http.StatusOK, "shared/invitation_index", gin.H{
				"Context":   c,
				"VersionID": sn.VersionID(),
				"CUser":     cu,
				"Games":     gs,
			})
		default:
			c.HTML(http.StatusOK, "shared/multi_games_index", gin.H{
				"Context":   c,
				"VersionID": sn.VersionID(),
				"CUser":     cu,
				"Games":     gs,
				"Status":    status,
			})
		}
	}
}

func (client *Client) JIndex(c *gin.Context) {
	client.Log.Debugf("Entering")
	defer client.Log.Debugf("Exiting")

	options := struct {
		ItemsPerPage int    `json:"itemsPerPage"`
		Forward      string `json:"forward"`
	}{}

	err := c.ShouldBind(&options)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	client.Log.Debugf("options: %#v", options)

	cu, err := client.User.Current(c)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	forward, err := datastore.DecodeCursor(options.Forward)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	status := ToStatus[c.Param("status")]
	q := datastore.
		NewQuery("Game").
		Filter("Status=", int(status)).
		Order("-UpdatedAt")

	cnt, err := client.DS.Count(c, q)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	client.Log.Debugf("cnt: %v", cnt)
	items := options.ItemsPerPage
	if options.ItemsPerPage == -1 {
		items = cnt
	}

	var es []*withID
	it := client.DS.Run(c, q.Start(forward))
	for i := 0; i < items; i++ {
		var gh Header
		_, err := it.Next(&gh)
		if err == iterator.Done {
			break
		}
		if err != nil {
			sn.JErr(c, err)
			return
		}
		es = append(es, &withID{&gh})
	}

	forward, err = it.Cursor()
	if err != nil {
		sn.JErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"gheaders":   es,
		"totalItems": cnt,
		"forward":    forward.String(),
		"cu":         cu,
	})
}

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

type dbState interface {
	DBState()
}

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

func (h *Header) undoKey(cu *user.User) string {
	return fmt.Sprintf("%s/uid-%d", h.Key, cu.ID())
}

func (h *Header) UndoKey(cu *user.User) string {
	return h.undoKey(cu)
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

func (client Client) JSONIndexAction(c *gin.Context) {
	client.Log.Debugf("Entering")
	defer client.Log.Debugf("Exiting")

	cu, err := client.User.Current(c)
	if err != nil {
		client.Log.Warningf(err.Error())
	}

	data, err := toGameTable(c, cu)
	if err != nil {
		client.Log.Errorf(err.Error())
		c.JSON(http.StatusOK, fmt.Sprintf("%v", err))
		return
	}
	c.JSON(http.StatusOK, data)
}

func toGameTable(c *gin.Context, cu *user.User) (*jGamesIndex, error) {
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
			Players:     h.PlayerLinks(cu),
			NumPlayers:  template.HTML(fmt.Sprintf("%d / %d", h.AcceptedPlayers(), h.NumPlayers)),
			Round:       h.Round,
			OptString:   template.HTML(h.OptString),
			Progress:    template.HTML(h.Progress),
			UpdatedAt:   h.UpdatedAt,
			LastUpdated: template.HTML(restful.LastUpdated(time.Time(h.UpdatedAt))),
			Public:      publicPrivate(g),
			Actions:     actionButtons(c, cu, g),
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

func ToGameTable(c *gin.Context, gs []Gamer, cnt int64, cu *user.User) (*jGamesIndex, error) {
	return toGameTable(withCount(withGamers(c, gs), cnt), cu)
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
		<div style="font-size:.7em">%s</div>`, h.Type.IDString(), h.ID(), h.Title, h.OptString))
}

func actionButtons(c *gin.Context, cu *user.User, g Gamer) template.HTML {
	log.Debugf("Entering")
	defer log.Debugf("Exiting")

	h := g.GetHeader()
	switch h.Status {
	case Running:
		t := h.Type.IDString()
		if g.GetHeader().IsCurrentPlayer(cu) {
			return template.HTML(fmt.Sprintf(`<a class="mybutton" href="/%s/game/show/%d">Play</a>`, t, h.ID()))
		} else {
			return template.HTML(fmt.Sprintf(`<a class="mybutton" href="/%s/game/show/%d">Show</a>`, t, h.ID()))
		}
	case Recruiting:
		t := h.Type.IDString()
		switch {
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
	</script>`, h.ID(), h.ID(), h.Stub(), h.ID(), h.ID(), h.ID(), h.ID(), h.ID()))
			} else {
				return template.HTML(fmt.Sprintf(`
				<form method="post" action="/%s/game/accept/%d">
					<input name="_method" type="hidden" value="PUT" />
					<input id="user_id" name="user[id]" type="hidden" value="%v">
					<input id="accept-%d" class="mybutton" type="submit" value="Accept" />
				</form>`, t, h.ID(), cu.ID(), h.ID()))
			}
		case g.CanDropout(cu):
			return template.HTML(fmt.Sprintf(`
				<form method="post" action="/%s/game/drop/%d">
					<input name="_method" type="hidden" value="PUT" />
					<input id="user_id" name="user[id]" type="hidden" value="%v">
					<input id="drop-%d" class="mybutton" type="submit" value="Drop" />
				</form>`, t, h.ID(), cu.ID(), h.ID()))
		default:
			return ""
		}
	default:
		return ""
	}
}

func (cl *Client) GamesIndex(c *gin.Context) {
	cl.Log.Debugf("Entering")
	defer cl.Log.Debugf("Exiting")

	obj := struct {
		Options struct {
			ItemsPerPage int `json:"itemsPerPage"`
		} `json:"options"`
		Forward string `json:"forward"`
		Status  string `json:"status"`
		Type    string `json:"type"`
	}{}

	err := c.ShouldBind(&obj)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	cl.Log.Debugf("obj: %#v", obj)

	cu, err := cl.User.Current(c)
	if err != nil {
		sn.JErr(c, err)
		return
	}
	cl.Log.Debugf("cu: %#v", cu)
	cl.Log.Debugf("err: %#v", err)

	forward, err := datastore.DecodeCursor(obj.Forward)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	cl.Log.Debugf("forward: %#v", forward)
	status := ToStatus[obj.Status]
	t := gType.ToType[obj.Type]
	q := datastore.
		NewQuery("Game").
		Filter("Status=", int(status)).
		Order("-UpdatedAt")

	if t != gType.All && t != gType.NoType {
		q = q.Filter("Type=", int(t))
	}

	cnt, err := cl.DS.Count(c, q)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	cl.Log.Debugf("cnt: %v", cnt)
	items := obj.Options.ItemsPerPage
	if obj.Options.ItemsPerPage == -1 {
		items = cnt
	}

	var es []*GHeader
	it := cl.DS.Run(c, q.Start(forward))
	for i := 0; i < items; i++ {
		var h Header
		k, err := it.Next(&h)
		if err == iterator.Done {
			break
		}
		if err != nil {
			sn.JErr(c, err)
			return
		}
		es = append(es, &GHeader{Key: k, Header: h})
	}

	forward, err = it.Cursor()
	if err != nil {
		sn.JErr(c, err)
		return
	}

	cl.Log.Debugf("forward: %#v", forward)
	cl.Log.Debugf("forward.String: %#v", forward.String())
	c.JSON(http.StatusOK, gin.H{
		"gheaders":   es,
		"totalItems": cnt,
		"forward":    forward.String(),
		"cu":         cu,
	})
}
