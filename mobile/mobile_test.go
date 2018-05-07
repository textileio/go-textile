package mobile_test

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"testing"

	"github.com/segmentio/ksuid"

	"github.com/textileio/textile-go/core"
	. "github.com/textileio/textile-go/mobile"
	util "github.com/textileio/textile-go/util/testing"

	libp2p "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

var wrapper *Wrapper
var hash string

var cusername = ksuid.New().String()
var cpassword = ksuid.New().String()
var cemail = ksuid.New().String() + "@textile.io"

func TestNewTextile(t *testing.T) {
	os.RemoveAll("testdata/.ipfs")
	var err error
	wrapper, err = NewNode("testdata/.ipfs", util.CentralApiURL)
	if err != nil {
		t.Errorf("create mobile node failed: %s", err)
	}
}

func TestWrapper_Start(t *testing.T) {
	err := wrapper.Start()
	if err != nil {
		t.Errorf("start mobile node failed: %s", err)
	}
}

func TestWrapper_StartAgain(t *testing.T) {
	err := wrapper.Start()
	if err != nil {
		t.Errorf("attempt to start a running node failed: %s", err)
	}
}

func TestWrapper_SignUpWithEmail(t *testing.T) {
	_, ref, err := util.CreateReferral(util.RefKey, 1)
	if err != nil {
		t.Errorf("create referral for signup failed: %s", err)
		return
	}
	if len(ref.RefCodes) == 0 {
		t.Error("create referral for signup got no codes")
		return
	}
	err = wrapper.SignUpWithEmail(cusername, cpassword, cemail, ref.RefCodes[0])
	if err != nil {
		t.Errorf("signup failed: %s", err)
		return
	}
}

func TestWrapper_SignIn(t *testing.T) {
	err := wrapper.SignIn(cusername, cpassword)
	if err != nil {
		t.Errorf("signin failed: %s", err)
		return
	}
}

func TestWrapper_IsSignedIn(t *testing.T) {
	if !wrapper.IsSignedIn() {
		t.Errorf("is signed in check failed should be true")
		return
	}
}

func TestWrapper_GetUsername(t *testing.T) {
	un, err := wrapper.GetUsername()
	if err != nil {
		t.Errorf("get username failed: %s", err)
		return
	}
	if un != cusername {
		t.Errorf("got bad username: %s", un)
	}
}

func TestWrapper_GetAccessToken(t *testing.T) {
	_, err := wrapper.GetAccessToken()
	if err != nil {
		t.Errorf("get access token failed: %s", err)
		return
	}
}

func TestWrapper_GetGatewayPassword(t *testing.T) {
	pwd := wrapper.GetGatewayPassword()
	if pwd == "" {
		t.Errorf("got bad gateway password: %s", pwd)
		return
	}
}

func TestWrapper_AddPhoto(t *testing.T) {
	mr, err := wrapper.AddPhoto("testdata/image.jpg", "testdata/thumb.jpg", "default")
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

func TestWrapper_SharePhoto(t *testing.T) {
	mr, err := wrapper.SharePhoto(hash, "beta")
	if err != nil {
		t.Errorf("share photo failed: %s", err)
		return
	}
	if len(mr.Boundary) == 0 {
		t.Errorf("share photo got bad hash")
	}
}

func TestWrapper_GetPhotos(t *testing.T) {
	res, err := wrapper.GetPhotos("", -1, "default")
	if err != nil {
		t.Errorf("get photos failed: %s", err)
		return
	}
	list := core.PhotoList{}
	json.Unmarshal([]byte(res), &list)
	if len(list.Items) == 0 {
		t.Errorf("get photos bad result")
	}
}

func TestWrapper_GetPhotosEmptyChannel(t *testing.T) {
	res, err := wrapper.GetPhotos("", -1, "empty")
	if err != nil {
		t.Errorf("get photos failed: %s", err)
		return
	}
	list := core.PhotoList{}
	json.Unmarshal([]byte(res), &list)
	if len(list.Items) != 0 {
		t.Errorf("get photos bad result")
	}
}

func TestWrapper_GetFileBase64(t *testing.T) {
	res, err := wrapper.GetFileBase64(hash + "/thumb")
	if err != nil {
		t.Errorf("get photo base64 string failed: %s", err)
		return
	}
	if len(res) == 0 {
		t.Errorf("get photo base64 string bad result")
	}
}

func TestWrapper_GetPeerID(t *testing.T) {
	_, err := wrapper.GetPeerID()
	if err != nil {
		t.Errorf("get peer id failed: %s", err)
	}
}

func TestWrapper_PairDesktop(t *testing.T) {
	_, pk, err := libp2p.GenerateKeyPair(libp2p.RSA, 4096)
	if err != nil {
		t.Errorf("create rsa keypair failed: %s", err)
	}
	pb, err := pk.Bytes()
	if err != nil {
		t.Errorf("get rsa keypair bytes: %s", err)
	}
	ps := base64.StdEncoding.EncodeToString(pb)

	_, err = wrapper.PairDesktop(ps)
	if err != nil {
		t.Errorf("pair desktop failed: %s", err)
	}
}

func TestWrapper_SignOut(t *testing.T) {
	err := wrapper.SignOut()
	if err != nil {
		t.Errorf("signout failed: %s", err)
		return
	}
}

func TestWrapper_IsSignedInAgain(t *testing.T) {
	if wrapper.IsSignedIn() {
		t.Errorf("is signed in check failed should be false")
		return
	}
}

func TestWrapper_Stop(t *testing.T) {
	err := wrapper.Stop()
	if err != nil {
		t.Errorf("stop mobile node failed: %s", err)
	}
}

func TestWrapper_StopAgain(t *testing.T) {
	err := wrapper.Stop()
	if err != nil {
		t.Errorf("stop mobile node again should not return error: %s", err)
	}
}

func Test_Teardown(t *testing.T) {
	os.RemoveAll(wrapper.RepoPath)
}
