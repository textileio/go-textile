package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/textileio/go-textile/pb"
)

// addContacts godoc
// @Summary Add to known contacts
// @Description Adds a contact by username or account address to known contacts.
// @Tags contacts
// @Accept application/json
// @Param address path string true "address"
// @Param contact body pb.Contact true "contact"
// @Success 204 {string} string "ok"
// @Failure 400 {string} string "Bad Request"
// @Router /contacts/{address} [put]
func (a *Api) addContacts(g *gin.Context) {
	var contact pb.Contact
	if err := pbUnmarshaler.Unmarshal(g.Request.Body, &contact); err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}
	if contact.Address == "" || len(contact.Peers) == 0 {
		g.String(http.StatusBadRequest, "invalid contact")
		return
	}
	if contact.Address != g.Param("address") {
		g.String(http.StatusBadRequest, "contact address mismatch")
		return
	}

	if err := a.Node.AddContact(&contact); err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	a.Node.FlushCafes()

	g.Status(http.StatusNoContent)
}

// lsContacts godoc
// @Summary List known contacts
// @Description Lists known contacts.
// @Tags contacts
// @Produce application/json
// @Success 200 {object} pb.ContactList "contacts"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /contacts [get]
func (a *Api) lsContacts(g *gin.Context) {
	pbJSON(g, http.StatusOK, a.Node.Contacts())
}

// getContacts godoc
// @Summary Get a known contact
// @Description Gets a known contact
// @Tags contacts
// @Produce application/json
// @Param address path string true "address"
// @Success 200 {object} pb.Contact "contact"
// @Failure 404 {string} string "Not Found"
// @Router /contacts/{address} [get]
func (a *Api) getContacts(g *gin.Context) {
	contact := a.Node.Contact(g.Param("address"))
	if contact == nil {
		g.String(http.StatusNotFound, "contact not found")
		return
	}

	pbJSON(g, http.StatusOK, contact)
}

// rmContacts godoc
// @Summary Remove a contact
// @Description Removes a known contact
// @Tags contacts
// @Param address path string true "address"
// @Success 204 {string} string "ok"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /contacts/{address} [delete]
func (a *Api) rmContacts(g *gin.Context) {
	address := g.Param("address")

	contact := a.Node.Contact(address)
	if contact == nil {
		g.String(http.StatusNotFound, "contact not found")
		return
	}

	if err := a.Node.RemoveContact(address); err != nil {
		a.abort500(g, err)
		return
	}

	a.Node.FlushCafes()

	g.Status(http.StatusNoContent)
}

// searchContacts godoc
// @Summary Search for contacts
// @Description Search for contacts known locally and on the network
// @Tags contacts
// @Produce application/json
// @Param X-Textile-Opts header string false "local: Whether to only search local contacts, remote: Whether to only search remote contacts, limit: Stops searching after limit results are found, wait: Stops searching after 'wait' seconds have elapsed (max 30s), username: search by username string, address: search by account address string, events: Whether to emit Server-Sent Events (SSEvent) or plain JSON" default(local="false",limit=5,wait=5,address=,username=,events="false")
// @Success 200 {object} pb.QueryResult "results stream"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /contacts/search [post]
func (a *Api) searchContacts(g *gin.Context) {
	opts, err := a.readOpts(g)
	if err != nil {
		a.abort500(g, err)
		return
	}

	localOnly, err := strconv.ParseBool(opts["local"])
	if err != nil {
		localOnly = false
	}
	remoteOnly, err := strconv.ParseBool(opts["remote"])
	if err != nil {
		remoteOnly = false
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
		Address: opts["address"],
		Name:    opts["name"],
	}
	options := &pb.QueryOptions{
		LocalOnly:  localOnly,
		RemoteOnly: remoteOnly,
		Limit:      int32(limit),
		Wait:       int32(wait),
	}

	resCh, errCh, cancel, err := a.Node.SearchContacts(query, options)
	if err != nil {
		g.String(http.StatusBadRequest, err.Error())
		return
	}

	handleSearchStream(g, resCh, errCh, cancel, opts["events"] == "true")
}
