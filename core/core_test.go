package core_test

import (
	"github.com/op/go-logging"
	"github.com/segmentio/ksuid"
	cmodels "github.com/textileio/textile-go/central/models"
	. "github.com/textileio/textile-go/core"
	util "github.com/textileio/textile-go/util/testing"
	"os"
	"testing"
	"time"
)

var node *TextileNode
var hash string

var centralReg = &cmodels.Registration{
	Username: ksuid.New().String(),
	Password: ksuid.New().String(),
	Identity: &cmodels.Identity{
		Type:  cmodels.EmailAddress,
		Value: ksuid.New().String() + "@textile.io",
	},
	Referral: "",
}

func TestNewNode(t *testing.T) {
	os.RemoveAll("testdata/.ipfs")
	config := NodeConfig{
		RepoPath:      "testdata/.ipfs",
		CentralApiURL: util.CentralApiURL,
		IsMobile:      false,
		LogLevel:      logging.DEBUG,
		LogFiles:      false,
	}
	var err error
	node, err = NewNode(config)
	if err != nil {
		t.Errorf("create node failed: %s", err)
	}
}

func TestTextileNode_StartWallet(t *testing.T) {
	online, err := node.StartWallet()
	if err != nil {
		t.Errorf("start node failed: %s", err)
	}
}

func TestTextileNode_StartAgain(t *testing.T) {
	err := node.Start()
	if err != ErrNodeRunning {
		t.Errorf("start node again reported wrong error: %s", err)
	}
}

func TestTextileNode_StartGarbageCollection(t *testing.T) {
	_, err := node.StartGarbageCollection()
	if err != nil {
		t.Errorf("start services failed: %s", err)
	}
}

func TestTextileNode_Online(t *testing.T) {
	if !node.Online() {
		t.Errorf("should report online")
	}
}

func TestTextileNode_SignUp(t *testing.T) {
	_, ref, err := util.CreateReferral(util.RefKey, 1)
	if err != nil {
		t.Errorf("create referral for signup failed: %s", err)
		return
	}
	if len(ref.RefCodes) == 0 {
		t.Error("create referral for signup got no codes")
		return
	}
	centralReg.Referral = ref.RefCodes[0]

	err = node.SignUp(centralReg)
	if err != nil {
		t.Errorf("signup failed: %s", err)
		return
	}
}

func TestTextileNode_SignIn(t *testing.T) {
	creds := &cmodels.Credentials{
		Username: centralReg.Username,
		Password: centralReg.Password,
	}
	err := node.SignIn(creds)
	if err != nil {
		t.Errorf("signin failed: %s", err)
		return
	}
}

func TestTextileNode_IsSignedIn(t *testing.T) {
	// TODO
}

func TestTextileNode_GetUsername(t *testing.T) {
	// TODO
}

func TestTextileNode_GetAccessToken(t *testing.T) {
	// TODO
}

func TestTextileNode_JoinRoom(t *testing.T) {
	da := node.Datastore.Albums().GetAlbumByName("default")
	if da == nil {
		t.Error("default album not found")
		return
	}
	go node.JoinRoom(da.Id, make(chan ThreadUpdate))
}

func TestTextileNode_LeaveRoom(t *testing.T) {
	da := node.Datastore.Albums().GetAlbumByName("default")
	if da == nil {
		t.Error("default album not found")
		return
	}
	time.Sleep(time.Second)
	node.LeaveRoom(da.Id)
	<-node.LeftRoomChs[da.Id]
}

func TestTextileNode_WaitForRoom(t *testing.T) {
	// TODO
}

func TestTextileNode_CreateAlbum(t *testing.T) {
	err := node.CreateAlbum("", "test")
	if err != nil {
		t.Errorf("create album failed: %s", err)
		return
	}
}

func TestTextileNode_AddPhoto(t *testing.T) {
	caption := "i am not a crook"
	mr, err := node.AddPhoto("testdata/image.jpg", "testdata/thumb.jpg", "default", caption)
	if err != nil {
		t.Errorf("add photo failed: %s", err)
		return
	}
	if len(mr.Boundary) == 0 {
		t.Errorf("add photo got bad hash")
	}
	hash = mr.Boundary
	err = os.Remove("testdata/" + mr.Boundary)
	if err != nil {
		t.Errorf("error unlinking test multipart file: %s", err)
	}
}

func TestTextileNode_SharePhoto(t *testing.T) {
	caption := "a day that will live on in infamy"
	mr, err := node.SharePhoto(hash, "test", caption)
	if err != nil {
		t.Errorf("share photo failed: %s", err)
		return
	}
	if len(mr.Boundary) == 0 {
		t.Errorf("share photo got bad hash")
	}
}

func TestTextileNode_GetHashRequest(t *testing.T) {
	// TODO
}

func TestTextileNode_GetPhotos(t *testing.T) {
	list := node.GetPhotos("", -1, "default")
	if len(list.Hashes) == 0 {
		t.Errorf("get photos bad result")
	}
}

func TestTextileNode_GetFile(t *testing.T) {
	res, err := node.GetFile(hash+"/thumb", nil)
	if err != nil {
		t.Errorf("get photo failed: %s", err)
		return
	}
	if len(res) == 0 {
		t.Errorf("get photo bad result")
	}
}

func TestTextileNode_GetFileBase64(t *testing.T) {
	// TODO
}

func TestTextileNode_GetMetadata(t *testing.T) {
	_, err := node.GetMetaData(hash, nil)
	if err != nil {
		t.Errorf("get metadata failed: %s", err)
		return
	}
}

func TestTextileNode_GetLastHash(t *testing.T) {
	_, err := node.GetLastHash(hash, nil)
	if err != nil {
		t.Errorf("get last hash failed: %s", err)
		return
	}
}

func TestTextileNode_LoadPhotoAndAlbum(t *testing.T) {
	// TODO
}

func TestTextileNode_UnmarshalPrivatePeerKey(t *testing.T) {
	_, err := node.UnmarshalPrivatePeerKey()
	if err != nil {
		t.Errorf("unmarshal private peer key failed: %s", err)
		return
	}
}

func TestTextileNode_GetPublicPeerKeyString(t *testing.T) {
	_, err := node.GetPublicPeerKeyString()
	if err != nil {
		t.Errorf("get peer public key as base 64 string failed: %s", err)
		return
	}
}

func TestTextileNode_PingPeer(t *testing.T) {
	// TODO
}

func TestTextileNode_SignOut(t *testing.T) {
	err := node.SignOut()
	if err != nil {
		t.Errorf("signout failed: %s", err)
		return
	}
}

func TestTextileNode_Stop(t *testing.T) {
	err := node.Stop()
	if err != nil {
		t.Errorf("stop node failed: %s", err)
	}
}

func TestTextileNode_OnlineAgain(t *testing.T) {
	if node.Online() {
		t.Errorf("should report offline")
	}
}

// test signin in stopped state, should re-connect to db
func TestTextileNode_SignInAgain(t *testing.T) {
	creds := &cmodels.Credentials{
		Username: centralReg.Username,
		Password: centralReg.Password,
	}
	err := node.SignIn(creds)
	if err != nil {
		t.Errorf("signin failed: %s", err)
		return
	}
}

func Test_Teardown(t *testing.T) {
	os.RemoveAll(node.RepoPath)
}
