package game

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/SlothNinja/color"
	"github.com/SlothNinja/rating"
	"github.com/SlothNinja/user"
	stats "github.com/SlothNinja/user-stats"
	"github.com/gin-gonic/gin"
)

func init() {
	gob.Register(NewPlayer())
}

type Comparison int

const (
	EqualTo     Comparison = 0
	LessThan    Comparison = -1
	GreaterThan Comparison = 1
)

const NoPlayerID int = -1

type Player struct {
	gamer           Gamer
	user            *user.User
	rating          *rating.CurrentRating
	stats           *stats.Stats
	IDF             int  `form:"idf"`
	PerformedAction bool `form:"performed-action"`
	Score           int  `form:"score"`
	Passed          bool `form:"passed"`
	ColorMapF       color.Colors
}
type Players []*Player

type jPlayer struct {
	User            *user.User            `json:"user"`
	Rating          *rating.CurrentRating `json:"rating"`
	IDF             int                   `json:"id"`
	PerformedAction bool                  `json:"performedAction"`
	Score           int                   `json:"score"`
	Passed          bool                  `json:"passed"`
	ColorMap        []string              `json:"colorMap"`
	//	Link            string                `json:"link"`
	//Gravatar string `json:"gravatar"`
}

func (p *Player) MarshalJSON() ([]byte, error) {
	j := &jPlayer{
		User:            p.user,
		Rating:          p.rating,
		IDF:             p.IDF,
		PerformedAction: p.PerformedAction,
		Score:           p.Score,
		Passed:          p.Passed,
		ColorMap:        p.gamer.(GetPlayerers).GetPlayerers().Colors().Strings(),
		//		Link:            string(p.Link()),
		// Gravatar: p.Gravatar(),
	}
	return json.Marshal(j)
}

type Playerer interface {
	ID() int
	Index() int
	User() *user.User
	Name() string
	Color() color.Color
	ColorMap() color.Colors
	// Init(Gamer)
	// Rating() *rating.CurrentRating
	Stats() *stats.Stats
}

func (p *Player) CompareByScore(p2 *Player) (c Comparison) {
	switch {
	case p.Score < p2.Score:
		c = LessThan
	case p.Score > p2.Score:
		c = GreaterThan
	default:
		c = EqualTo
	}
	return
}

type Playerers []Playerer

func (p *Player) Game() Gamer {
	return p.gamer
}

func (p *Player) SetGame(gamer Gamer) {
	p.gamer = gamer
}

func NewPlayer() (p *Player) {
	p = new(Player)
	return
}

func (p *Player) ID() int {
	return p.IDF
}

func (p *Player) SetID(id int) {
	p.IDF = id
}

func (p *Player) ColorMap() color.Colors {
	return p.ColorMapF
}

func (p *Player) SetColorMap(colors color.Colors) {
	p.ColorMapF = colors
}

func (p *Player) Equal(p2 Playerer) bool {
	return p2 != nil && p.ID() == p2.ID()
}

func (p *Player) NotEqual(p2 Playerer) bool {
	return !p.Equal(p2)
}

func (p *Player) User() *user.User {
	if p.user == nil {
		p.user = p.gamer.User(p.ID())
	}
	return p.user
}

func (h *Header) UserIDFor(p Playerer) (id int64) {
	if l, pid := len(h.UserIDS), p.ID(); pid >= 0 && pid < l {
		id = h.UserIDS[p.ID()]
	}
	return
}

func (h *Header) NameFor(p Playerer) (n string) {
	if p != nil {
		n = h.NameByPID(p.ID())
	}
	return
}

func (h *Header) NameByPID(pid int) (n string) {
	if l := len(h.UserNames); pid >= 0 && pid < l {
		n = h.UserNames[pid]
	}
	return
}

func (h *Header) NameByUID(uid int64) (n string) {
	var index int = NotFound
	for i := range h.UserIDS {
		if uid == h.UserIDS[i] {
			index = i
			break
		}
	}

	if index != NotFound {
		n = h.NameByPID(index)
	}
	return
}

func (h *Header) NameByUSID(sid string) (n string) {
	var index int = NotFound
	for i := range h.UserIDS {
		if sid == h.UserSIDS[i] {
			index = i
			break
		}
	}

	if index != NotFound {
		n = h.NameByPID(index)
	}
	return
}

func (h *Header) EmailFor(p Playerer) (em string) {
	if l, pid := len(h.UserEmails), p.ID(); pid >= 0 && pid < l {
		em = h.UserEmails[p.ID()]
	}
	return
}

// func (p *Player) Rating() *rating.CurrentRating {
// 	if p.rating != nil {
// 		return p.rating
// 	}
//
// 	c := p.User().CTX()
// 	dsClient, err := datastore.NewClient(c, "")
// 	if err != nil {
// 		log.Warningf("unable to open datastore client")
// 	}
//
// 	client := rating.NewClient(dsClient)
// 	p.rating, _ = client.For(c, p.User(), p.Game().GetHeader().Type)
// 	return p.rating
// }

func (p *Player) Stats() *stats.Stats {
	if p.stats != nil {
		return p.stats
	}
	p.stats = p.gamer.Stat(p.ID())
	return p.stats
}

// Name provides the name of the player.
// TODO: Deprecated in favor of NameFor.
func (p *Player) Name() (s string) {
	if p != nil && p.User() != nil {
		s = p.User().Name
	}
	return
}

// Index provides the index in players for the player.
// TODO: Deprecated in favor of IndexFor
func (p *Player) Index() (index int) {
	return IndexFor(p, p.Game().(GetPlayerers).GetPlayerers())
}

// IndexFor returns the index for the player in players, if present.
// Returns NotFound, the player not in players.
func IndexFor(p Playerer, ps Playerers) (index int) {
	index = NotFound
	for i, p2 := range ps {
		if p.ID() == p2.ID() {
			index = i
			break
		}
	}
	return
}

// NotFound indicates a value (e.g., player) was not found in the collection.
const NotFound = -1

func (p *Player) Color() color.Color {
	if p == nil {
		return color.None
	}
	colorMap := p.gamer.DefaultColorMap()
	if cu := p.gamer.CurrentUser(); cu != nil {
		if player := p.gamer.PlayererByUserID(cu.ID()); player != nil {
			colorMap = player.ColorMap()
		}
	}
	return colorMap[p.ID()]
}

func (ps Playerers) Colors() color.Colors {
	cs := make(color.Colors, len(ps))
	for i, p := range ps {
		cs[i] = p.Color()
	}
	return cs
}

var textColors = map[color.Color]color.Color{
	color.Yellow: color.Black,
	color.Purple: color.White,
	color.Green:  color.Yellow,
	color.White:  color.Black,
	color.Black:  color.White,
}

func (p *Player) TextColor() (c color.Color) {
	var ok bool
	if c, ok = textColors[p.Color()]; !ok {
		c = color.Black
	}
	return
}

// A bit of a misnomer
// Returns whether the current user is an admin
//func (p *Player) IsAdmin() bool {
//	if p == nil {
//		return false
//	}
//	return p.gamer.CurrentUser().IsAdmin(p.gamer.CTX())
//}

// A bit of a misnomer
// Returns whether the current user is the same as the player's user
func (p *Player) IsCurrentUser(c *gin.Context) bool {
	if p == nil {
		return false
	}
	cu, err := user.CurrentFrom(c)
	if err != nil {
		return false
	}
	return p.User().Equal(cu)
}

// A bit of a misnomer
// Returns whether the current user is the same as the player's user or
// whether the current user is an admin.
//func (p *Player) IsCurrentUserOrAdmin(ctx context.Context) bool {
//	if p == nil {
//		return false
//	}
//	return p.IsCurrentUser() || user.IsAdmin(ctx)
//}

func (p *Player) IsCurrentPlayer() bool {
	for _, player := range p.gamer.CurrentPlayerers() {
		if player != nil && p.Equal(player) {
			return true
		}
	}
	return false
}

func (p *Player) IsWinner() (b bool) {
	for _, p2 := range p.gamer.Winnerers() {
		if b = p2 != nil && p.Equal(p2); b {
			break
		}
	}
	return
}

//func (p *Player) Link() template.HTML {
//	u := p.User()
//	cp := p.IsCurrentPlayer()
//	me := u.IsCurrent(p.gamer.CTX())
//	w := p.IsWinner()
//
//	result := fmt.Sprintf(`<a href="/user/show/%d" >%s</a>`, u.ID, u.Name)
//	h := p.gamer.GetHeader()
//	switch h.Status {
//	case Running:
//		switch {
//		case cp && me:
//			result = fmt.Sprintf(`<a href="/user/show/%d" class="current-player me">%s</a>`, u.ID, u.Name)
//		case cp:
//			result = fmt.Sprintf(`<a href="/user/show/%d" class="current-player">%s</a>`, u.ID, u.Name)
//		}
//	case Completed:
//		switch {
//		case w && me:
//			result = fmt.Sprintf(`<a href="/user/show/%d" class="winner me">%s</a>`, u.ID, u.Name)
//		case w:
//			result = fmt.Sprintf(`<a href="/user/show/%d" class="winner">%s</a>`, u.ID, u.Name)
//		}
//	}
//	return template.HTML(result)
//}

func (p *Player) Init(g Gamer) {
	p.SetGame(g)
}

func (p *Player) Gravatar() string {
	return fmt.Sprintf(`<a href="/user/show/%d" ><img src=%q alt="Gravatar" class="%s-border" /> </a>`,
		p.User().ID, p.User().Gravatar("80"), p.Color())
}

func (h *Header) GravatarFor(p Playerer) template.HTML {
	return template.HTML(fmt.Sprintf(`<a href=%q ><img src=%q alt="Gravatar" class="%s-border" /> </a>`,
		h.UserPathFor(p), user.GravatarURL(h.EmailFor(p), "monsterid", "80"), p.Color()))
}

func (h *Header) UserPathFor(p Playerer) template.HTML {
	return user.PathFor(h.UserIDFor(p))
}
