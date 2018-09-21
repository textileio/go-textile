package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
)

func ListNotifications(c *ishell.Context) {
	notifs := core.Node.GetNotifications("", -1)
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
		var username string
		if notif.ActorUsername != "" {
			username = notif.ActorUsername
		} else {
			username = notif.ActorId
		}
		c.Println(yellow(fmt.Sprintf("%s: #%s: %s %s.", notif.Id, notif.Subject, username, body)))
	}
}

func ReadNotification(c *ishell.Context) {
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

func ReadAllNotifications(c *ishell.Context) {
	if err := core.Node.ReadAllNotifications(); err != nil {
		c.Err(err)
		return
	}

	yellow := color.New(color.FgHiYellow).SprintFunc()
	c.Println(yellow("ok, marked all as read"))
}

func AcceptThreadInviteViaNotification(c *ishell.Context) {
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
