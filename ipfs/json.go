package ipfs

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/proto"
	icid "github.com/ipfs/go-cid"
	files "github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	format "github.com/ipfs/go-ipld-format"
	uio "github.com/ipfs/go-unixfs/io"
	ipb "github.com/ipfs/go-unixfs/pb"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/path"
)

// AddJSON adds JSON as a node and returns its top-level content ID
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

// JSONAtPath returns the JSON representation of the node at the given path
func JSONAtPath(node *core.IpfsNode, pth string) ([]byte, error) {
	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(node.Context(), CatTimeout)
	defer cancel()

	j, err := toJSONValue(ctx, api, pth)
	if err != nil {
		return nil, err
	}

	return json.Marshal(j)
}

func addToDir(ctx context.Context, node *core.IpfsNode, api iface.CoreAPI, k string, v interface{}, dir uio.Directory) error {
	switch v.(type) {
	case nil:
		return nil
	case string, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		return addDataToDir(ctx, api, k, bytes.NewReader(data), dir)
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
			return addDirToDir(ctx, node, api, k, d, dir)
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
				return addDirToDir(ctx, node, api, k, d, dir)
			} else {
				return err
			}
		}
	}
}

func addDataToDir(ctx context.Context, api iface.CoreAPI, k string, v io.Reader, dir uio.Directory) error {
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

func toJSONValue(ctx context.Context, api iface.CoreAPI, pth string) (interface{}, error) {
	ipth := path.New(pth)
	nd, err := api.Object().Get(ctx, ipth)
	if err != nil {
		return nil, err
	}
	r, err := api.Object().Data(ctx, ipth)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	pbdata := new(ipb.Data)
	err = proto.Unmarshal(data, pbdata)
	if err != nil {
		return nil, err
	}
	switch *pbdata.Type {
	case ipb.Data_File:
		var v interface{}
		err = json.Unmarshal(pbdata.Data, &v)
		if err != nil {
			return nil, err
		}
		return v, nil
	case ipb.Data_Directory:
		if !isListNode(nd) {
			m := make(map[string]interface{})
			for _, v := range nd.Links() {
				m[v.Name], err = toJSONValue(ctx, api, v.Cid.String())
				if err != nil {
					return nil, err
				}
			}
			return m, nil
		} else {
			s := make([]interface{}, len(nd.Links()))
			for i, l := range nd.Links() {
				s[i], err = toJSONValue(ctx, api, l.Cid.String())
				if err != nil {
					return nil, err
				}
			}
			return s, nil
		}
	default:
		return nil, nil
	}
}

func isListNode(n format.Node) bool {
	for i, l := range n.Links() {
		k, err := strconv.Atoi(l.Name)
		if err != nil || k != i {
			return false
		}
	}
	return true
}
