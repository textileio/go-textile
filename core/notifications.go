package core

import (
	"errors"
	"fmt"

	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"

	"github.com/textileio/textile-go/pb"
)

// Notifications lists notifications
func (t *Textile) Notifications(offset string, limit int) *pb.NotificationList {
	list := t.datastore.Notifications().List(offset, limit)
	for i, note := range list.Items {
		list.Items[i] = t.NotificationView(note)
	}
	return list
}

// NotificationView returns a notification with expanded view info
func (t *Textile) NotificationView(note *pb.Notification) *pb.Notification {
	switch note.Type {
	case pb.Notification_INVITE_RECEIVED:
		invite := t.InviteView(t.datastore.Invites().Get(note.Block))
		if invite != nil {
			note.User = invite.Inviter
		}
	default:
		note.User = t.User(note.Actor)
	}
	return note
}

// CountUnreadNotifications counts unread notifications
func (t *Textile) CountUnreadNotifications() int {
	return t.datastore.Notifications().CountUnread()
}

// ReadNotification marks a notification as read
func (t *Textile) ReadNotification(id string) error {
	return t.datastore.Notifications().Read(id)
}

// ReadAllNotifications marks all notification as read
func (t *Textile) ReadAllNotifications() error {
	return t.datastore.Notifications().ReadAll()
}

// AcceptInviteViaNotification uses an invite notification to accept an invite to a thread
func (t *Textile) AcceptInviteViaNotification(id string) (mh.Multihash, error) {
	notification := t.datastore.Notifications().Get(id)
	if notification == nil {
		return nil, errors.New(fmt.Sprintf("could not find notification: %s", id))
	}
	if notification.Type != pb.Notification_INVITE_RECEIVED {
		return nil, errors.New(fmt.Sprintf("notification type is not invite"))
	}

	hash, err := t.AcceptInvite(notification.Block)
	if err != nil {
		return nil, err
	}

	return hash, t.datastore.Notifications().Delete(id)
}

// IgnoreInviteViaNotification uses an invite notification to ignore an invite to a thread
func (t *Textile) IgnoreInviteViaNotification(id string) error {
	notification := t.datastore.Notifications().Get(id)
	if notification == nil {
		return errors.New(fmt.Sprintf("could not find notification: %s", id))
	}
	if notification.Type != pb.Notification_INVITE_RECEIVED {
		return errors.New(fmt.Sprintf("notification type is not invite"))
	}

	if err := t.IgnoreInvite(notification.Block); err != nil {
		return err
	}

	return t.datastore.Notifications().Delete(id)
}
