package cmd

import (
	"fmt"
	"github.com/textileio/textile-go/core"
	"gopkg.in/abiosoft/ishell.v2"
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

	// add node commands
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

	// add all commands w/ shell counterparts
	for _, c := range Cmds() {
		if c.Shell() != nil {
			shell.AddCmd(c.Shell())
		}
	}

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
			Name: "rm",
			Help: "remove a cafe",
			Func: cafeDeregister,
		})
		cafeCmd.AddCmd(&ishell.Cmd{
			Name: "ls",
			Help: "list active cafes",
			Func: cafeList,
		})
		cafeCmd.AddCmd(&ishell.Cmd{
			Name: "check-messages",
			Help: "check messages from active cafes",
			Func: cafeCheckMessages,
		})
		shell.AddCmd(cafeCmd)
	}
	{
		profileCmd := &ishell.Cmd{
			Name:     "profile",
			Help:     "manage profile",
			LongHelp: "Resolve other profiles, get and publish own profile.",
		}
		profileCmd.AddCmd(&ishell.Cmd{
			Name: "publish",
			Help: "publish profile",
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
			Name: "avatar",
			Help: "set profile avatar",
			Func: setAvatar,
		})
		profileCmd.AddCmd(&ishell.Cmd{
			Name: "username",
			Help: "set profile username",
			Func: setUsername,
		})
		profileCmd.AddCmd(&ishell.Cmd{
			Name: "subs",
			Help: "list profile subs",
			Func: getSubs,
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
		//photoCmd.AddCmd(&ishell.Cmd{
		//	Name: "add",
		//	Help: "add a new photo",
		//	Func: addPhoto,
		//})
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
			Help: "show key for a photo (and meta data)",
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
