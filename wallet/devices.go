package wallet

import (
	"errors"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/wallet/thread"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
	"time"
)

// Devices lists all devices
func (w *Wallet) Devices() []repo.Device {
	return w.datastore.Devices().List("")
}

// AddDevice creates an invite for every current and future thread
func (w *Wallet) AddDevice(name string, pk libp2pc.PubKey) error {
	if !w.IsOnline() {
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
	if err := w.datastore.Devices().Add(deviceModel); err != nil {
		return err
	}
	log.Infof("added device '%s'", name)

	// invite device to existing threads
	for _, thrd := range w.threads {
		if _, err := thrd.AddInvite(pk); err != nil {
			return err
		}
	}

	// notify listeners
	w.sendUpdate(Update{Id: deviceModel.Id, Name: deviceModel.Name, Type: DeviceAdded})

	// send notification
	id, err := w.GetId()
	if err != nil {
		return err
	}
	notification := &repo.Notification{
		Id:            ksuid.New().String(),
		Date:          time.Now(),
		ActorId:       id,
		ActorUsername: "You",
		Subject:       deviceModel.Name,
		SubjectId:     deviceModel.Id,
		Type:          repo.DeviceAddedNotification,
		Body:          "paired with a new device",
	}
	return w.sendNotification(notification)
}

// InviteDevices sends a thread invite to all devices
func (w *Wallet) InviteDevices(thrd *thread.Thread) error {
	for _, device := range w.Devices() {
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
func (w *Wallet) RemoveDevice(id string) error {
	if !w.IsOnline() {
		return ErrOffline
	}

	// delete db record
	device := w.datastore.Devices().Get(id)
	if device == nil {
		return errors.New("device not found")
	}
	if err := w.datastore.Devices().Delete(id); err != nil {
		return err
	}

	// delete notifications
	if err := w.datastore.Notifications().DeleteBySubjectId(device.Id); err != nil {
		return err
	}

	log.Infof("removed device '%s'", id)

	// notify listeners
	w.sendUpdate(Update{Id: device.Id, Name: device.Name, Type: DeviceRemoved})

	return nil
}
