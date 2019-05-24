package cmd

import (
	"fmt"
	"github.com/textileio/go-textile/common"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/gateway"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
	"github.com/textileio/go-textile/util"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"
)

// Start the daemon against the user repository
func Daemon(repoPath string, pinCode string, docs bool, debug bool) error {
	var err error
	node, err = core.NewTextile(core.RunConfig{
		PinCode:  pinCode,
		RepoPath: repoPath,
		Debug:    debug,
	})
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("create node failed: %s", err))
	}

	gateway.Host = &gateway.Gateway{
		Node: node,
	}

	if err := startNode(docs); err != nil {
		return fmt.Errorf(fmt.Sprintf("start node failed: %s", err))
	}
	printSplash()

	// Shutdown gracefully if an SIGINT was received
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Interrupted")
	fmt.Printf("Shutting down...")
	if err := stopNode(); err != nil && err != core.ErrStopped {
		fmt.Println(err.Error())
	} else {
		fmt.Print("done\n")
	}
	os.Exit(1)
	return nil
}


// Output the instance environment for the daemon command
func printSplash() {
	pid, err := node.PeerId()
	if err != nil {
		log.Fatalf("get peer id failed: %s", err)
	}
	fmt.Println(Grey("go-textile version: " + common.GitSummary))
	fmt.Println(Grey("Repo version: ") + Grey(repo.Repover))
	fmt.Println(Grey("Repo path: ") + Grey(node.RepoPath()))
	fmt.Println(Grey("API address: ") + Grey(node.ApiAddr()))
	fmt.Println(Grey("Gateway address: ") + Grey(gateway.Host.Addr()))
	if node.CafeApiAddr() != "" {
		fmt.Println(Grey("Cafe address: ") + Grey(node.CafeApiAddr()))
	}
	fmt.Println(Grey("System version: ") + Grey(runtime.GOARCH+"/"+runtime.GOOS))
	fmt.Println(Grey("Golang version: ") + Grey(runtime.Version()))
	fmt.Println(Grey("PeerID:  ") + Green(pid.Pretty()))
	fmt.Println(Grey("Account: ") + Cyan(node.Account().Address()))
}


// Start the node, the API, and the Gateway
// And subsribe to updates of the wallet, thread, and notifications
func startNode(serveDocs bool) error {
	listener := node.ThreadUpdateListener()

	if err := node.Start(); err != nil {
		return err
	}

	// subscribe to wallet updates
	go func() {
		for {
			select {
			case update, ok := <-node.UpdateCh():
				if !ok {
					return
				}
				switch update.Type {
				case pb.WalletUpdate_THREAD_ADDED:
					break
				case pb.WalletUpdate_THREAD_REMOVED:
					break
				case pb.WalletUpdate_ACCOUNT_PEER_ADDED:
					break
				case pb.WalletUpdate_ACCOUNT_PEER_REMOVED:
					break
				}
			}
		}
	}()

	// subscribe to thread updates
	go func() {
		for {
			select {
			case value, ok := <-listener.Ch:
				if !ok {
					return
				}
				if update, ok := value.(*pb.FeedItem); ok {
					thrd := update.Thread[len(update.Thread)-8:]

					btype, err := core.FeedItemType(update)
					if err != nil {
						log.Error(err.Error())
						continue
					}

					payload, err := core.GetFeedItemPayload(update)
					if err != nil {
						log.Error(err.Error())
						continue
					}
					user := payload.GetUser()
					date := payload.GetDate()

					var txt string
					txt += time.Unix(0, util.ProtoNanos(date)).Format(time.RFC822)
					txt += "  "

					if user != nil {
						var name string
						if user.Name != "" {
							name = user.Name
						} else {
							if len(user.Address) >= 7 {
								name = user.Address[:7]
							} else {
								name = user.Address
							}
						}
						txt += name + " "
					}
					txt += "added "

					msg := Grey(txt) + Green(btype.String()) + Grey(" update to "+thrd)
					fmt.Println(msg)
				}
			}
		}
	}()

	// subscribe to notifications
	go func() {
		for {
			select {
			case note, ok := <-node.NotificationCh():
				if !ok {
					return
				}

				date := util.ProtoTime(note.Date).Format(time.RFC822)
				var subject string
				if len(note.Subject) >= 7 {
					subject = note.Subject[len(note.Subject)-7:]
				}

				msg := Grey(date+"  "+note.User.Name+" ") + Cyan(note.Body) +
					Grey(" "+subject)
				fmt.Println(msg)
			}
		}
	}()

	// start apis
	node.StartApi(node.Config().Addresses.API, serveDocs)
	gateway.Host.Start(node.Config().Addresses.Gateway)

	// start profiling api
	go func() {
		writeHeapDump("/debug/write-heap-dump/")
		freeOSMemory("/debug/free-os-memory/")
		mutexFractionOption("/debug/pprof-mutex/")
		if err := http.ListenAndServe(node.Config().Addresses.Profiling, http.DefaultServeMux); err != nil {
			log.Errorf("error staring profile listener: %s", err)
		}
	}()

	// Wait concurrently here until the node comes online
	// that is to say, until the online channel opens
	<-node.OnlineCh()

	// Textile is now online, continue
	return nil
}

// Stop the api, then the gateway, then the node, then if possible, the channels
// If a former fails, do not continue with the latter
func stopNode() error {
	if err := node.StopApi(); err != nil {
		return err
	}
	if err := gateway.Host.Stop(); err != nil {
		return err
	}
	if err := node.Stop(); err != nil {
		return err
	}

	node.CloseChns()
	return nil
}

// mutexFractionOption allows to set runtime.SetMutexProfileFraction via HTTP
// using POST request with parameter 'fraction'.
func mutexFractionOption(path string) {
	http.DefaultServeMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		asfr := r.Form.Get("fraction")
		if len(asfr) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fr, err := strconv.Atoi(asfr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		log.Infof("Setting MutexProfileFraction to %d", fr)
		runtime.SetMutexProfileFraction(fr)
	})
}

// writeHeapDump writes a description of the heap and the objects in
// it to the given file descriptor. (used here for debugging)
func writeHeapDump(path string) {
	http.DefaultServeMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		log.Infof("Writing heap dump")
		f, err := os.Create("heapdump")
		if err != nil {
			return
		}
		debug.WriteHeapDump(f.Fd())
	})
}

// freeOSMemory forces a garbage collection followed by an
// attempt to return as much memory to the operating system
// as possible. (used here for debugging)
func freeOSMemory(path string) {
	http.DefaultServeMux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		log.Infof("Freeing OS memory")
		debug.FreeOSMemory()
	})
}
