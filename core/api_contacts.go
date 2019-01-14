package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *api) lsContacts(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	var contacts []ContactInfo

	threadId := opts["thread"]
	if threadId == "default" {
		threadId = a.node.config.Threads.Defaults.ID
	}
	if threadId != "" {
		thrd := a.node.Thread(threadId)
		if thrd == nil {
			g.String(http.StatusNotFound, ErrThreadNotFound.Error())
			return
		}
		contacts = make([]ContactInfo, 0)
		for _, p := range thrd.Peers() {
			contact := a.node.Contact(p.Id)
			if contact != nil {
				contacts = append(contacts, *contact)
			}
		}
	} else {
		contacts, err = a.node.Contacts()
	}

	g.JSON(http.StatusOK, contacts)
}

func (a *api) getContacts(g *gin.Context) {
	id := g.Param("id")

	info := a.node.Contact(id)
	if info == nil {
		g.String(http.StatusNotFound, "contact not found")
		return
	}

	g.JSON(http.StatusOK, info)
}
