package ipfs

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"strconv"
	"strings"

	icid "github.com/ipfs/go-cid"
	files "github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	uio "github.com/ipfs/go-unixfs/io"
	iface "github.com/ipfs/interface-go-ipfs-core"
)

// AddJSON adds a JSON string as a Directory
func AddJSON(node *core.IpfsNode, input string) (*icid.Cid, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(node.Context(), PinTimeout)
	defer cancel()

	dir := uio.NewDirectory(api.Dag())
	var m map[string]interface{}
	err = json.Unmarshal([]byte(input), &m)
	if err == nil {
		for k, v := range m {
			err = addToDir(ctx, node, api, k, v, dir)
			if err != nil {
				return nil, err
			}
		}
	} else {
		var s []interface{}
		err = json.Unmarshal([]byte(input), &s)
		if err == nil {
			for i, v := range s {
				err = addToDir(ctx, node, api, strconv.Itoa(i), v, dir)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	n, err := dir.GetNode()
	if err != nil {
		return nil, err
	}

	defer node.Blockstore.PinLock().Unlock()
	err = node.Pinning.Pin(ctx, n, false)
	if err != nil && !strings.Contains(err.Error(), "already pinned recursively") {
		return nil, err
	}
	err = node.Pinning.Flush()
	if err != nil {
		return nil, err
	}

	cid := n.Cid()
	return &cid, nil
}

func addToDir(ctx context.Context, node *core.IpfsNode, api iface.CoreAPI, k string, v interface{}, dir uio.Directory) error {
	buf := new(bytes.Buffer)
	switch v.(type) {
	case nil:
		return nil
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		err := binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			return err
		}
	case string:
		_, err := buf.Write([]byte(v.(string)))
		if err != nil {
			return err
		}
	case error:
		err := binary.Write(buf, binary.LittleEndian, v.(error).Error())
		if err != nil {
			return err
		}
	default:
		d := uio.NewDirectory(api.Dag())
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		var m map[string]interface{}
		err = json.Unmarshal(b, &m)
		if err == nil {
			for k, v := range m {
				err = addToDir(ctx, node, api, k, v, d)
				if err != nil {
					return err
				}
			}
			err = addDirToDir(ctx, node, api, k, d, dir)
			if err != nil {
				return err
			}
		} else {
			var s []interface{}
			err = json.Unmarshal(b, &s)
			if err == nil {
				for i, v := range s {
					err = addToDir(ctx, node, api, strconv.Itoa(i), v, d)
					if err != nil {
						return err
					}
				}
				err = addDirToDir(ctx, node, api, k, d, dir)
				if err != nil {
					return err
				}
			}
		}
	}

	return addDataToDir(ctx, api, k, buf, dir)
}

func addDataToDir(ctx context.Context, api iface.CoreAPI, k string, v *bytes.Buffer, dir uio.Directory) error {
	if v.Len() == 0 {
		return nil
	}
	pth, err := api.Unixfs().Add(ctx, files.NewReaderFile(v))
	if err != nil {
		return err
	}
	n, err := api.Dag().Get(ctx, pth.Cid())
	if err != nil {
		return err
	}
	return dir.AddChild(ctx, k, n)
}

func addDirToDir(ctx context.Context, node *core.IpfsNode, api iface.CoreAPI, k string, v uio.Directory, dir uio.Directory) error {
	nd, err := v.GetNode()
	if err != nil {
		return err
	}

	defer node.Blockstore.PinLock().Unlock()
	err = node.Pinning.Pin(ctx, nd, false)
	if err != nil && !strings.Contains(err.Error(), "already pinned recursively") {
		return err
	}
	err = node.Pinning.Flush()
	if err != nil {
		return err
	}

	n, err := api.Dag().Get(ctx, nd.Cid())
	if err != nil {
		return err
	}
	return dir.AddChild(ctx, k, n)
}
