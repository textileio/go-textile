package core

import (
	"errors"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/thread"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"time"
)

// Devices lists all devices
func (t *Textile) Devices() []repo.Device {
	return t.datastore.Devices().List("")
}

// AddDevice creates an invite for every current and future thread
func (t *Textile) AddDevice(name string, pk libp2pc.PubKey) error {
	if !t.IsOnline() {
		return ErrOffline
	}

	// index a new device
	pkb, err := pk.Bytes()
	if err != nil {
		return err
	}
	deviceModel := &repo.Device{
		Id:   libp2pc.ConfigEncodeKey(pkb),
		Name: name,
	}
	if err := t.datastore.Devices().Add(deviceModel); err != nil {
		return err
	}
	log.Infof("added device '%s'", name)

	// invite device to existing threads
	for _, thrd := range t.threads {
		if _, err := thrd.AddInvite(pk); err != nil {
			return err
		}
	}

	// notify listeners
	t.sendUpdate(Update{Id: deviceModel.Id, Name: deviceModel.Name, Type: DeviceAdded})

	// send notification
	id, err := t.GetID()
	if err != nil {
		return err
	}
	notification := &repo.Notification{
		Id:            ksuid.New().String(),
		Date:          time.Now(),
		ActorId:       id.Pretty(),
		ActorUsername: "You",
		Subject:       deviceModel.Name,
		SubjectId:     deviceModel.Id,
		Type:          repo.DeviceAddedNotification,
		Body:          "paired with a new device",
	}
	return t.sendNotification(notification)
}

// InviteDevices sends a thread invite to all devices
func (t *Textile) InviteDevices(thrd *thread.Thread) error {
	for _, device := range t.Devices() {
		dpkb, err := libp2pc.ConfigDecodeKey(device.Id)
		if err != nil {
			return err
		}
		dpk, err := libp2pc.UnmarshalPublicKey(dpkb)
		if err != nil {
			return err
		}
		if _, err := thrd.AddInvite(dpk); err != nil {
			return err
		}
	}
	return nil
}

// RemoveDevice removes a device
func (t *Textile) RemoveDevice(id string) error {
	if !t.IsOnline() {
		return ErrOffline
	}

	// delete db record
	device := t.datastore.Devices().Get(id)
	if device == nil {
		return errors.New("device not found")
	}
	if err := t.datastore.Devices().Delete(id); err != nil {
		return err
	}

	// delete notifications
	if err := t.datastore.Notifications().DeleteBySubjectId(device.Id); err != nil {
		return err
	}

	log.Infof("removed device '%s'", id)

	// notify listeners
	t.sendUpdate(Update{Id: device.Id, Name: device.Name, Type: DeviceRemoved})

	return nil
}
