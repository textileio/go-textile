package core

import (
	"fmt"
	"strconv"

	"github.com/textileio/go-textile/ipfs"
	"github.com/textileio/go-textile/pb"
)

func (t *Textile) Files(offset string, limit int, threadId string) (*pb.FilesList, error) {
	var query string
	if threadId != "" {
		if t.Thread(threadId) == nil {
			return nil, ErrThreadNotFound
		}
		query = fmt.Sprintf("threadId='%s' and type=%d", threadId, pb.Block_FILES)
	} else {
		query = fmt.Sprintf("type=%d", pb.Block_FILES)
	}

	list := make([]*pb.Files, 0)

	blocks := t.Blocks(offset, limit, query)
	for _, block := range blocks.Items {
		file, err := t.file(block, feedItemOpts{annotations: true})
		if err != nil {
			return nil, err
		}
		list = append(list, file)
	}

	return &pb.FilesList{Items: list}, nil
}

func (t *Textile) File(blockId string) (*pb.Files, error) {
	block, err := t.Block(blockId)
	if err != nil {
		return nil, err
	}

	return t.file(block, feedItemOpts{annotations: true})
}

func (t *Textile) fileAtData(data string) ([]*pb.File, error) {
	links, err := ipfs.LinksAtPath(t.node, data)
	if err != nil {
		return nil, err
	}

	files := make([]*pb.File, len(links))

	for _, index := range links {
		node, err := ipfs.NodeAtLink(t.node, index)
		if err != nil {
			return nil, err
		}

		i, err := strconv.Atoi(index.Name)
		if err != nil {
			return nil, err
		}

		f := &pb.File{Index: int32(i)}
		if looksLikeFileNode(node) {
			file, err := t.fileIndexForPair(node)
			if err != nil {
				return nil, err
			}
			f.File = file

		} else {
			f.Links = make(map[string]*pb.FileIndex)
			for _, link := range node.Links() {
				pair, err := ipfs.NodeAtLink(t.node, link)
				if err != nil {
					return nil, err
				}
				file, err := t.fileIndexForPair(pair)
				if err != nil {
					return nil, err
				}
				if file != nil {
					f.Links[link.Name] = file
				}
			}
		}

		files[i] = f
	}

	return files, nil
}

func (t *Textile) file(block *pb.Block, opts feedItemOpts) (*pb.Files, error) {
	if block.Type != pb.Block_FILES {
		return nil, ErrBlockWrongType
	}

	files, err := t.fileAtData(block.Data)
	if err != nil {
		return nil, err
	}

	item := &pb.Files{
		Block:   block.Id,
		Data:    block.Data,
		Date:    block.Date,
		User:    t.PeerUser(block.Author),
		Caption: block.Body,
		Files:   files,
		Threads: t.fileThreads(block.Data),
	}

	if opts.annotations {
		comments, err := t.Comments(block.Id)
		if err != nil {
			return nil, err
		}
		item.Comments = comments.Items

		likes, err := t.Likes(block.Id)
		if err != nil {
			return nil, err
		}
		item.Likes = likes.Items
	} else {
		item.Comments = opts.comments
		item.Likes = opts.likes
	}

	return item, nil
}

// fileThreads lists threads that have blocks which link a file
// @todo This should be a distinct db query, if it's even needed?
func (t *Textile) fileThreads(data string) []string {
	unique := make([]string, 0)
	threads := make(map[string]struct{})

	for _, b := range t.datastore.Blocks().List("", -1, "data='"+data+"'").Items {
		if _, ok := threads[b.Thread]; !ok {
			threads[b.Thread] = struct{}{}
			unique = append(unique, b.Thread)
		}
	}
	return unique
}
