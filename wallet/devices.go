package wallet

import (
	"errors"
	"fmt"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/repo"
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
		Id:       ksuid.New().String(),
		Date:     time.Now(),
		ActorId:  id,
		TargetId: deviceModel.Id,
		Type:     repo.DeviceAddedNotification,
		Body:     fmt.Sprintf("You are now paired with a new device named %s.", deviceModel.Name),
	}
	return w.sendNotification(notification)
}

// RemoveDevice removes a device
func (w *Wallet) RemoveDevice(id string) error {
	if !w.IsOnline() {
		return ErrOffline
	}

	device := w.datastore.Devices().Get(id)
	if device == nil {
		return errors.New("device not found")
	}
	if err := w.datastore.Devices().Delete(id); err != nil {
		return err
	}
	log.Infof("removed device '%s'", id)

	// TODO: uninvite?

	// notify listeners
	w.sendUpdate(Update{Id: device.Id, Name: device.Name, Type: DeviceRemoved})

	return nil
}
