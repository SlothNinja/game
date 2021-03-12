package game

import (
	"bytes"

	"github.com/SlothNinja/restful"
	"github.com/SlothNinja/send"
	gType "github.com/SlothNinja/type"
	"github.com/SlothNinja/user"
	"github.com/gin-gonic/gin"
	"github.com/mailjet/mailjet-apiv3-go"
)

const (
	sender  = "webmaster@slothninja.com"
	subject = "SlothNinja Games: Daily Turn Notifications"
)

type inf struct {
	GameID int64
	Type   gType.Type
	Title  string
}

type infs []*inf
type notifications map[int64]infs

func (client Client) DailyNotifications(c *gin.Context) {
	client.Log.Debugf("Entering")
	defer client.Log.Debugf("Exiting")

	gs := GamersFrom(c)

	notifications := make(notifications, 0)
	for _, g := range gs {
		h := g.GetHeader()
		gameInfo := &inf{GameID: h.ID(), Type: h.Type, Title: h.Title}
		for _, index := range h.CPUserIndices {
			uid := h.UserIDS[index]
			notifications[uid] = append(notifications[uid], gameInfo)
		}
	}

	msg := mailjet.InfoMessagesV31{
		From: &mailjet.RecipientV31{
			Email: "webmaster@slothninja.com",
			Name:  "Webmaster",
		},
		Subject: subject,
	}
	tmpl := restful.TemplatesFrom(c)["shared/daily_notification"]
	buf := new(bytes.Buffer)

	for uid, gameInfos := range notifications {
		m := msg
		u := user.New(uid)

		err := client.DS.Get(c, u.Key, u)
		if err != nil {
			client.Log.Errorf("get user error: %s", err.Error())
			buf.Reset()
			continue
		}

		if !u.EmailReminders {
			continue
		}

		err = tmpl.Execute(buf, gin.H{
			"Info": gameInfos,
			"User": u,
		})
		if err != nil {
			client.Log.Errorf("template execution for %s generated error: %s", u.Name, err.Error())
			buf.Reset()
			continue
		}

		m.HTMLPart = buf.String()
		m.To = &mailjet.RecipientsV31{
			mailjet.RecipientV31{
				Email: u.Email,
				Name:  u.Name,
			},
		}

		_, err = send.Messages(c, m)
		if err != nil {
			client.Log.Errorf("enqueuing email message: %#v geneerated error: %s", m, err.Error())
			buf.Reset()
			continue
		}

		// Reset buffer for next message
		buf.Reset()
	}
}
