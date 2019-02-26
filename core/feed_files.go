package core

import (
	"fmt"
	"strconv"

	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/pb"
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

func (t *Textile) fileAtTarget(target string) ([]*pb.File, error) {
	links, err := ipfs.LinksAtPath(t.node, target)
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

		info := &pb.File{Index: int32(i)}
		if looksLikeFileNode(node) {
			file, err := t.fileIndexForPair(node)
			if err != nil {
				return nil, err
			}
			info.File = file

		} else {
			info.Links = make(map[string]*pb.FileIndex)
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
					info.Links[link.Name] = file
				}
			}
		}

		files[i] = info
	}

	return files, nil
}

func (t *Textile) file(block *pb.Block, opts feedItemOpts) (*pb.Files, error) {
	if block.Type != pb.Block_FILES {
		return nil, ErrBlockWrongType
	}

	threads := make([]string, 0)
	threads = t.fileThreads(block.Target)

	files, err := t.fileAtTarget(block.Target)
	if err != nil {
		return nil, err
	}

	item := &pb.Files{
		Block:   block.Id,
		Target:  block.Target,
		Date:    block.Date,
		User:    t.User(block.Author),
		Caption: block.Body,
		Files:   files,
		Threads: threads,
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

// fileThreads lists threads that have blocks which target a file
func (t *Textile) fileThreads(target string) []string {
	var unique []string

	blocks := t.datastore.Blocks().List("", -1, "target='"+target+"'")
outer:
	for _, b := range blocks.Items {
		for _, f := range unique {
			if f == b.Thread {
				break outer
			}
		}
		unique = append(unique, b.Thread)
	}

	return unique
}
