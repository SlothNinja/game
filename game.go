package game

import (
	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/color"
	"github.com/SlothNinja/log"
	"github.com/SlothNinja/sn"
	gtype "github.com/SlothNinja/type"
	"github.com/SlothNinja/user"
	"github.com/gin-gonic/gin"
)

type Gamers []Gamer
type Gamer interface {
	PhaseName() string
	FromParams(*gin.Context, *user.User, gtype.Type) error
	ColorMapFor(*user.User) color.Map
	headerer
}

type GetPlayerers interface {
	GetPlayerers() Playerers
}

type hasUpdate interface {
	Update(*gin.Context) (string, ActionType, error)
}

func GamesRoot(c *gin.Context) *datastore.Key {
	return datastore.NameKey("Games", "root", nil)
}

func (h *Header) GetAcceptDialog() bool {
	return h.Private()
}

func (h *Header) RandomTurnOrder() {
	ps := h.gamer.(GetPlayerers).GetPlayerers()
	for i := 0; i < h.NumPlayers; i++ {
		ri := sn.MyRand.Intn(h.NumPlayers)
		ps[i], ps[ri] = ps[ri], ps[i]
	}
	h.SetCurrentPlayerers(ps[0])

	h.OrderIDS = make(UserIndices, len(ps))
	for i, p := range ps {
		h.OrderIDS[i] = p.ID()
	}
}

// Returns (true, nil) if game should be started
func (h *Header) Accept(c *gin.Context, u *user.User) (start bool, err error) {
	log.Debugf("Entering")
	defer log.Debugf("Entering")

	if err = h.validateAccept(c, u); err != nil {
		return
	}

	h.AddUser(u)
	if len(h.Users) == h.NumPlayers {
		start = true
	}
	return
}

func (h *Header) validateAccept(c *gin.Context, u *user.User) error {
	switch {
	case len(h.UserIDS) >= h.NumPlayers:
		return sn.NewVError("Game already has the maximum number of players.")
	case h.HasUser(u):
		return sn.NewVError("%s has already accepted this invitation.", u.Name)
	case h.Password != "" && c.PostForm("password") != h.Password:
		return sn.NewVError("%s provided incorrect password for Game #%d: %s.", u.Name, h.ID, h.Title)
	}
	return nil
}

func (h *Header) Drop(u *user.User) (err error) {
	if err = h.validateDrop(u); err != nil {
		return
	}

	h.RemoveUser(u)
	return
}

func (h *Header) validateDrop(u *user.User) (err error) {
	switch {
	case h.Status != Recruiting:
		err = sn.NewVError("Game is no longer recruiting, thus %s can't drop.", u.Name)
	case !h.HasUser(u):
		err = sn.NewVError("%s has not joined this game, thus %s can't drop.", u.Name, u.Name)
	}
	return
}

// func RequireCurrentPlayerOrAdmin() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		log.Debugf("Entering")
// 		defer log.Debugf("Exiting")
//
// 		g := GamerFrom(c)
// 		if g == nil {
// 			log.Warningf("Missing Gamer")
// 			c.Abort()
// 			return
// 		}
//
// 		if !g.GetHeader().CUserIsCPlayerOrAdmin(c) {
// 			log.Warningf("Current User is Not Current Player or Admin")
// 			c.Abort()
// 			return
// 		}
// 	}
// }
