package cmd

import (
	"errors"
	"fmt"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	"os"
)

var shell *ishell.Shell

func RunShell(startNode func() error, stopNode func() error) {
	shell = ishell.New()
	shell.SetHomeHistoryPath(".ishell_history")

	// handle interrupt
	shell.Interrupt(func(c *ishell.Context, count int, input string) {
		if count == 1 {
			shell.Println("input Ctrl-C once more to exit")
			return
		}
		shell.Println("interrupted")
		shell.Printf("shutting down...")
		if err := stopNode(); err != nil && err != core.ErrStopped {
			c.Err(err)
		} else {
			shell.Printf("done\n")
		}
		os.Exit(1)
	})

	// add interactive commands
	shell.AddCmd(&ishell.Cmd{
		Name: "start",
		Help: "start the node",
		Func: func(c *ishell.Context) {
			if core.Node.Started() {
				c.Println("already started")
				return
			}
			if err := startNode(); err != nil {
				c.Println(fmt.Errorf("start node failed: %s", err))
				return
			}
			c.Println("ok, started")
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "stop",
		Help: "stop the node",
		Func: func(c *ishell.Context) {
			if !core.Node.Started() {
				c.Println("already stopped")
				return
			}
			if err := stopNode(); err != nil {
				c.Println(fmt.Errorf("stop node failed: %s", err))
				return
			}
			c.Println("ok, stopped")
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "id",
		Help: "show address and peer info",
		Func: showId,
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "ping",
		Help: "ping another peer",
		Func: func(c *ishell.Context) {
			if !core.Node.Online() {
				c.Println("not online yet")
				return
			}
			if len(c.Args) == 0 {
				c.Err(errors.New("missing peer id"))
				return
			}
			pid, err := peer.IDB58Decode(c.Args[0])
			if err != nil {
				c.Println(fmt.Errorf("bad peer id: %s", err))
				return
			}
			status, err := core.Node.Ping(pid)
			if err != nil {
				c.Println(fmt.Errorf("ping failed: %s", err))
				return
			}
			c.Println(status)
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "fetch-messages",
		Help: "fetch messages from registered cafes",
		Func: func(c *ishell.Context) {
			if !core.Node.Online() {
				c.Println("not online yet")
				return
			}
			if err := core.Node.FetchCafeMessages(); err != nil {
				c.Println(fmt.Errorf("fetch messages failed: %s", err))
				return
			}
			c.Println("ok, fetching")
		},
	})
	{
		cafeCmd := &ishell.Cmd{
			Name:     "cafe",
			Help:     "manage cafe sessions",
			LongHelp: "Manage cafe sessions.",
		}
		cafeCmd.AddCmd(&ishell.Cmd{
			Name: "add",
			Help: "add a new cafe",
			Func: cafeRegister,
		})
		cafeCmd.AddCmd(&ishell.Cmd{
			Name: "ls",
			Help: "list active cafes",
			Func: cafeList,
		})
		cafeCmd.AddCmd(&ishell.Cmd{
			Name: "rm",
			Help: "remove a cafe",
			Func: cafeDeregister,
		})
		shell.AddCmd(cafeCmd)
	}
	{
		profileCmd := &ishell.Cmd{
			Name:     "profile",
			Help:     "manage profile",
			LongHelp: "Resolve other profiles, get and publish local profile.",
		}
		profileCmd.AddCmd(&ishell.Cmd{
			Name: "publish",
			Help: "publish local profile",
			Func: publishProfile,
		})
		profileCmd.AddCmd(&ishell.Cmd{
			Name: "resolve",
			Help: "resolve profiles",
			Func: resolveProfile,
		})
		profileCmd.AddCmd(&ishell.Cmd{
			Name: "get",
			Help: "get peer profiles",
			Func: getProfile,
		})
		profileCmd.AddCmd(&ishell.Cmd{
			Name: "set-avatar",
			Help: "set local profile avatar",
			Func: setAvatarId,
		})
		shell.AddCmd(profileCmd)
	}
	{
		swarmCmd := &ishell.Cmd{
			Name:     "swarm",
			Help:     "same as ipfs swarm",
			LongHelp: "Inspect IPFS swarm peers.",
		}
		swarmCmd.AddCmd(&ishell.Cmd{
			Name: "peers",
			Help: "show connected peers (same as `ipfs swarm peers`)",
			Func: swarmPeers,
		})
		swarmCmd.AddCmd(&ishell.Cmd{
			Name: "ping",
			Help: "ping a peer (same as `ipfs ping`)",
			Func: swarmPing,
		})
		swarmCmd.AddCmd(&ishell.Cmd{
			Name: "connect",
			Help: "connect to a peer (same as `ipfs swarm connect`)",
			Func: swarmConnect,
		})
		shell.AddCmd(swarmCmd)
	}
	{
		photoCmd := &ishell.Cmd{
			Name:     "photo",
			Help:     "manage photos",
			LongHelp: "Add, list, and get info about photos.",
		}
		photoCmd.AddCmd(&ishell.Cmd{
			Name: "add",
			Help: "add a new photo",
			Func: addPhoto,
		})
		photoCmd.AddCmd(&ishell.Cmd{
			Name: "share",
			Help: "share a photo to a different thread",
			Func: sharePhoto,
		})
		photoCmd.AddCmd(&ishell.Cmd{
			Name: "get",
			Help: "save a photo to a local file",
			Func: getPhoto,
		})
		photoCmd.AddCmd(&ishell.Cmd{
			Name: "key",
			Help: "decrypt and print the key for a photo",
			Func: getPhotoKey,
		})
		photoCmd.AddCmd(&ishell.Cmd{
			Name: "meta",
			Help: "get photo metadata",
			Func: getPhotoMetadata,
		})
		photoCmd.AddCmd(&ishell.Cmd{
			Name: "ls",
			Help: "list photos from a thread",
			Func: listPhotos,
		})
		photoCmd.AddCmd(&ishell.Cmd{
			Name: "comment",
			Help: "comment on a photo (terminate input w/ ';'",
			Func: addPhotoComment,
		})
		photoCmd.AddCmd(&ishell.Cmd{
			Name: "like",
			Help: "like a photo",
			Func: addPhotoLike,
		})
		photoCmd.AddCmd(&ishell.Cmd{
			Name: "comments",
			Help: "list photo comments",
			Func: listPhotoComments,
		})
		photoCmd.AddCmd(&ishell.Cmd{
			Name: "likes",
			Help: "list photo likes",
			Func: listPhotoLikes,
		})
		shell.AddCmd(photoCmd)
	}
	{
		threadCmd := &ishell.Cmd{
			Name:     "thread",
			Help:     "manage threads",
			LongHelp: "Add, remove, list, invite to, and get info about textile threads.",
		}
		threadCmd.AddCmd(&ishell.Cmd{
			Name: "add",
			Help: "add a new thread",
			Func: addThread,
		})
		threadCmd.AddCmd(&ishell.Cmd{
			Name: "rm",
			Help: "remove a thread by name",
			Func: removeThread,
		})
		threadCmd.AddCmd(&ishell.Cmd{
			Name: "ls",
			Help: "list threads",
			Func: listThreads,
		})
		threadCmd.AddCmd(&ishell.Cmd{
			Name: "blocks",
			Help: "list blocks",
			Func: listThreadBlocks,
		})
		threadCmd.AddCmd(&ishell.Cmd{
			Name: "head",
			Help: "show current HEAD",
			Func: getThreadHead,
		})
		threadCmd.AddCmd(&ishell.Cmd{
			Name: "ignore",
			Help: "ignore a block",
			Func: ignoreBlock,
		})
		threadCmd.AddCmd(&ishell.Cmd{
			Name: "peers",
			Help: "list peers",
			Func: listThreadPeers,
		})
		threadCmd.AddCmd(&ishell.Cmd{
			Name: "invite",
			Help: "invite a peer to a thread",
			Func: addThreadInvite,
		})
		threadCmd.AddCmd(&ishell.Cmd{
			Name: "accept",
			Help: "accept a thread invite",
			Func: acceptThreadInvite,
		})
		threadCmd.AddCmd(&ishell.Cmd{
			Name: "invite-external",
			Help: "create an external invite link",
			Func: addExternalThreadInvite,
		})
		threadCmd.AddCmd(&ishell.Cmd{
			Name: "accept-external",
			Help: "accept an external thread invite",
			Func: acceptExternalThreadInvite,
		})
		shell.AddCmd(threadCmd)
	}
	{
		deviceCmd := &ishell.Cmd{
			Name:     "device",
			Help:     "manage connected devices",
			LongHelp: "Add, remove, and list connected devices.",
		}
		deviceCmd.AddCmd(&ishell.Cmd{
			Name: "add",
			Help: "add a new device",
			Func: addDevice,
		})
		deviceCmd.AddCmd(&ishell.Cmd{
			Name: "rm",
			Help: "remove a device by name",
			Func: removeDevice,
		})
		deviceCmd.AddCmd(&ishell.Cmd{
			Name: "ls",
			Help: "list devices",
			Func: listDevices,
		})
		shell.AddCmd(deviceCmd)
	}
	{
		notificationCmd := &ishell.Cmd{
			Name:     "notification",
			Help:     "manage notifications",
			LongHelp: "List and read notifications.",
		}
		notificationCmd.AddCmd(&ishell.Cmd{
			Name: "read",
			Help: "mark a notification as read",
			Func: readNotification,
		})
		notificationCmd.AddCmd(&ishell.Cmd{
			Name: "readall",
			Help: "mark all notifications as read",
			Func: readAllNotifications,
		})
		notificationCmd.AddCmd(&ishell.Cmd{
			Name: "ls",
			Help: "list notifications",
			Func: listNotifications,
		})
		notificationCmd.AddCmd(&ishell.Cmd{
			Name: "accept",
			Help: "accept an invite via notification",
			Func: acceptThreadInviteViaNotification,
		})
		shell.AddCmd(notificationCmd)
	}

	shell.Run()
}
