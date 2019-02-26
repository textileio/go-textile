package core

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/textileio/textile-go/pb"
)

// addContacts godoc
// @Summary Add to known contacts
// @Description Adds a contact to list of known local contacts. A common workflow is to pipe
// @Description /contact/search results into this endpoint, just be sure you know what the results
// @Description of the search are before adding!
// @Tags contacts
// @Accept application/json
// @Produce application/json
// @Param id path string true "contact id"
// @Param contact body pb.Contact true "contact"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "Bad Request"
// @Router /contacts/{id} [put]
func (a *api) addContacts(g *gin.Context) {
	var contact *pb.Contact
	if err := pbUnmarshaler.Unmarshal(g.Request.Body, contact); err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	if contact == nil {
		g.String(http.StatusBadRequest, "missing contact")
		return
	}
	if contact.Id == "" || contact.Address == "" {
		g.String(http.StatusBadRequest, "invalid contact")
		return
	}
	if contact.Id != g.Param("id") {
		g.String(http.StatusBadRequest, "contact id mismatch")
		return
	}

	if err := a.node.AddContact(contact); err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	g.String(http.StatusOK, "ok")
}

// lsContacts godoc
// @Summary List known contacts
// @Description Lists all, or thread-based, contacts.
// @Tags contacts
// @Produce application/json
// @Param X-Textile-Opts header string false "thread: Thread ID (omit for all known contacts)" default(thread=)
// @Success 200 {object} pb.ContactList "contacts"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /contacts [get]
func (a *api) lsContacts(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	contacts := &pb.ContactList{Items: make([]*pb.Contact, 0)}

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

		for _, p := range thrd.Peers() {
			contact := a.node.Contact(p.Id)
			if contact != nil {
				contacts.Items = append(contacts.Items, contact)
			}
		}

	} else {
		contacts = a.node.Contacts()
	}

	pbJSON(g, http.StatusOK, contacts)
}

// getContacts godoc
// @Summary Get a known contact
// @Description Gets information about a known contact
// @Tags contacts
// @Produce application/json
// @Param id path string true "contact id"
// @Success 200 {object} pb.Contact "contact"
// @Failure 404 {string} string "Not Found"
// @Router /contacts/{id} [get]
func (a *api) getContacts(g *gin.Context) {
	id := g.Param("id")

	info := a.node.Contact(id)
	if info == nil {
		g.String(http.StatusNotFound, "contact not found")
		return
	}

	pbJSON(g, http.StatusOK, info)
}

// rmContacts godoc
// @Summary Remove a contact
// @Description Removes a known contact
// @Tags contacts
// @Produce text/plain
// @Param id path string true "contact id"
// @Success 200 {string} string "ok"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /contacts/{id} [delete]
func (a *api) rmContacts(g *gin.Context) {
	id := g.Param("id")

	info := a.node.Contact(id)
	if info == nil {
		g.String(http.StatusNotFound, "contact not found")
		return
	}

	if err := a.node.RemoveContact(id); err != nil {
		a.abort500(g, err)
		return
	}

	g.String(http.StatusOK, "ok")
}

// searchContacts godoc
// @Summary Search for contacts
// @Description Search for contacts known locally and on the network
// @Tags contacts
// @Produce application/json
// @Param X-Textile-Opts header string false "local: Whether to only search local contacts, limit: Stops searching after limit results are found, wait: Stops searching after 'wait' seconds have elapsed (max 10s), username: search by username string, peer: search by peer id string, address: search by account address string, events: Whether to emit Server-Sent Events (SSEvent) or plain JSON" default(local="false",limit=5,wait=5,peer=,address=,username=,events="false")
// @Success 200 {object} pb.QueryResult "results stream"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /contacts/search [post]
func (a *api) searchContacts(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	local, err := strconv.ParseBool(opts["local"])
	if err != nil {
		local = false
	}
	limit, err := strconv.Atoi(opts["limit"])
	if err != nil {
		limit = 5
	}
	wait, err := strconv.Atoi(opts["wait"])
	if err != nil {
		wait = 5
	}

	query := &pb.ContactQuery{
		Id:       opts["peer"],
		Address:  opts["address"],
		Username: opts["username"],
	}
	options := &pb.QueryOptions{
		Local: local,
		Limit: int32(limit),
		Wait:  int32(wait),
	}

	resCh, errCh, cancel, err := a.node.SearchContacts(query, options)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	handleSearchStream(g, resCh, errCh, cancel, opts["events"] == "true")
}
