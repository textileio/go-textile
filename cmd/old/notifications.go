package old

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
)

func listNotifications(c *ishell.Context) {
	notifs := core.Node.Notifications("", -1)
	unread := core.Node.CountUnreadNotifications()
	if len(notifs) == 0 {
		c.Println("no notifications found")
	} else {
		c.Println(fmt.Sprintf("found %d notifications, %d unread", len(notifs), unread))
	}

	yellow := color.New(color.FgHiYellow).SprintFunc()
	for _, notif := range notifs {
		body := notif.Body
		if !notif.Read {
			body += " (unread)"
		}
		username := core.Node.ContactUsername(notif.ActorId)
		c.Println(yellow(fmt.Sprintf("%s: #%s: %s %s.", notif.Id, notif.Subject, username, body)))
	}
}

func readNotification(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing notification id"))
		return
	}
	id := c.Args[0]

	if err := core.Node.ReadNotification(id); err != nil {
		c.Err(err)
		return
	}

	yellow := color.New(color.FgHiYellow).SprintFunc()
	c.Println(yellow("ok, marked as read"))
}

func readAllNotifications(c *ishell.Context) {
	if err := core.Node.ReadAllNotifications(); err != nil {
		c.Err(err)
		return
	}

	yellow := color.New(color.FgHiYellow).SprintFunc()
	c.Println(yellow("ok, marked all as read"))
}

func acceptThreadInviteViaNotification(c *ishell.Context) {
	if len(c.Args) == 0 {
		c.Err(errors.New("missing notification id"))
		return
	}
	id := c.Args[0]

	if _, err := core.Node.AcceptThreadInviteViaNotification(id); err != nil {
		c.Err(err)
		return
	}

	green := color.New(color.FgHiGreen).SprintFunc()
	c.Println(green("ok, accepted via notification"))
}
