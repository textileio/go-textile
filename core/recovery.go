package core

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/textileio/textile-go/broadcast"
	"github.com/textileio/textile-go/pb"
)

// FindThreadBackups searches the network for backups
func (t *Textile) FindThreadBackups(query *pb.ThreadBackupQuery, options *pb.QueryOptions) (<-chan *pb.QueryResult, <-chan error, *broadcast.Broadcaster, error) {
	payload, err := proto.Marshal(query)
	if err != nil {
		return nil, nil, nil, err
	}

	options.Filter = pb.QueryOptions_NO_FILTER

	resCh, errCh, cancel := t.search(&pb.Query{
		Type:    pb.QueryType_THREAD_BACKUPS,
		Options: options,
		Payload: &any.Any{
			TypeUrl: "/ThreadBackupQuery",
			Value:   payload,
		},
	})

	// transform results by decrypting
	tresCh := make(chan *pb.QueryResult)
	terrCh := make(chan error)
	go func() {
		for {
			select {
			case res, ok := <-resCh:
				if !ok {
					close(tresCh)
					return
				}

				backup := new(pb.CafeClientThread)
				if err := ptypes.UnmarshalAny(res.Value, backup); err != nil {
					terrCh <- err
					break
				}

				plaintext, err := t.account.Decrypt(backup.Ciphertext)
				if err != nil {
					terrCh <- err
					break
				}

				//cthrd := new(pb.CafeThread)
				//if err := proto.Unmarshal(plaintext, cthrd); err != nil {
				//	terrCh <- err
				//	break
				//}

				res.Value = &any.Any{
					TypeUrl: "/CafeThread",
					Value:   plaintext,
				}

				tresCh <- res

			case err := <-errCh:
				terrCh <- err
			}
		}
	}()

	return tresCh, terrCh, cancel, nil
}

//// ApplyThreadBackup dencrypts and adds a thread from a backup
//func (t *Textile) ApplyThreadBackup(backup *pb.CafeClientThread) error {
//	plaintext, err := t.account.Decrypt(backup.Ciphertext)
//	if err != nil {
//		return err
//	}
//
//
//}
