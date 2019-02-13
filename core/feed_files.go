package core

import (
	"fmt"
	"strconv"

	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
)

func (t *Textile) Files(offset string, limit int, threadId string) ([]*pb.FeedFiles, error) {
	var query string
	if threadId != "" {
		if t.Thread(threadId) == nil {
			return nil, ErrThreadNotFound
		}
		query = fmt.Sprintf("threadId='%s' and type=%d", threadId, repo.FilesBlock)
	} else {
		query = fmt.Sprintf("type=%d", repo.FilesBlock)
	}

	list := make([]*pb.FeedFiles, 0)

	blocks := t.Blocks(offset, limit, query)
	for _, block := range blocks {
		file, err := t.feedFile(&block, true)
		if err != nil {
			return nil, err
		}
		list = append(list, file)
	}

	return list, nil
}

func (t *Textile) FeedFile(blockId string) (*pb.FeedFiles, error) {
	block, err := t.Block(blockId)
	if err != nil {
		return nil, err
	}

	return t.feedFile(block, true)
}

func (t *Textile) fileAtTarget(target string) ([]*pb.FeedFile, error) {
	links, err := ipfs.LinksAtPath(t.node, target)
	if err != nil {
		return nil, err
	}

	files := make([]*pb.FeedFile, len(links))

	for _, index := range links {
		node, err := ipfs.NodeAtLink(t.node, index)
		if err != nil {
			return nil, err
		}

		i, err := strconv.Atoi(index.Name)
		if err != nil {
			return nil, err
		}

		info := &pb.FeedFile{Index: int32(i)}
		if looksLikeFileNode(node) {
			file, err := t.fileForPair(node)
			if err != nil {
				return nil, err
			}
			info.File = file

		} else {
			info.Links = &pb.Directory{Files: make(map[string]*pb.File)}
			for _, link := range node.Links() {
				pair, err := ipfs.NodeAtLink(t.node, link)
				if err != nil {
					return nil, err
				}
				file, err := t.fileForPair(pair)
				if err != nil {
					return nil, err
				}
				if file != nil {
					info.Links.Files[link.Name] = file
				}
			}
		}

		files[i] = info
	}

	return files, nil
}

func (t *Textile) feedFile(block *repo.Block, annotated bool) (*pb.FeedFiles, error) {
	if block.Type != repo.FilesBlock {
		return nil, ErrBlockWrongType
	}

	threads := make([]string, 0)
	threads = t.fileThreads(block.Target)

	username, avatar := t.ContactDisplayInfo(block.AuthorId)

	files, err := t.fileAtTarget(block.Target)
	if err != nil {
		return nil, err
	}

	date, err := ptypes.TimestampProto(block.Date)
	if err != nil {
		return nil, err
	}

	info := &pb.FeedFiles{
		Block:    block.Id,
		Target:   block.Target,
		Date:     date,
		Author:   block.AuthorId,
		Username: username,
		Avatar:   avatar,
		Caption:  block.Body,
		Files:    files,
		Threads:  threads,
	}

	if annotated {
		comments, err := t.Comments(block.Id)
		if err != nil {
			return nil, err
		}
		info.Comments = comments

		likes, err := t.Likes(block.Id)
		if err != nil {
			return nil, err
		}
		info.Likes = likes
	}

	return info, nil
}

// fileThreads lists threads that have blocks which target a file
func (t *Textile) fileThreads(target string) []string {
	var unique []string

	blocks := t.datastore.Blocks().List("", -1, "target='"+target+"'")
outer:
	for _, b := range blocks {
		for _, f := range unique {
			if f == b.ThreadId {
				break outer
			}
		}
		unique = append(unique, b.ThreadId)
	}

	return unique
}
