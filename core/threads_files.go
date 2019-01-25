package core

import (
	"fmt"
	"strconv"
	"time"

	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/repo"
)

type ThreadFileInfo struct {
	Index int        `json:"index"`
	File  *repo.File `json:"file,omitempty"`
	Links Directory  `json:"links,omitempty"`
}

type ThreadFilesInfo struct {
	Block    string              `json:"block"`
	Target   string              `json:"target"`
	Date     time.Time           `json:"date"`
	AuthorId string              `json:"author_id"`
	Username string              `json:"username,omitempty"`
	Avatar   string              `json:"avatar,omitempty"`
	Caption  string              `json:"caption,omitempty"`
	Files    []ThreadFileInfo    `json:"files"`
	Comments []ThreadCommentInfo `json:"comments"`
	Likes    []ThreadLikeInfo    `json:"likes"`
	Threads  []string            `json:"threads"`
}

func (t *Textile) ThreadFiles(offset string, limit int, threadId string) ([]ThreadFilesInfo, error) {
	var query string
	if threadId != "" {
		if t.Thread(threadId) == nil {
			return nil, ErrThreadNotFound
		}
		query = fmt.Sprintf("threadId='%s' and type=%d", threadId, repo.FilesBlock)
	} else {
		query = fmt.Sprintf("type=%d", repo.FilesBlock)
	}

	list := make([]ThreadFilesInfo, 0)

	blocks := t.Blocks(offset, limit, query)
	for _, block := range blocks {
		file, err := t.threadFile(block)
		if err != nil {
			return nil, err
		}
		list = append(list, *file)
	}

	return list, nil
}

func (t *Textile) ThreadFile(blockId string) (*ThreadFilesInfo, error) {
	block, err := t.Block(blockId)
	if err != nil {
		return nil, err
	}

	return t.threadFile(*block)
}

func (t *Textile) fileAtTarget(target string) ([]ThreadFileInfo, error) {
	links, err := ipfs.LinksAtPath(t.node, target)
	if err != nil {
		return nil, err
	}

	files := make([]ThreadFileInfo, len(links))

	for _, index := range links {
		node, err := ipfs.NodeAtLink(t.node, index)
		if err != nil {
			return nil, err
		}

		i, err := strconv.Atoi(index.Name)
		if err != nil {
			return nil, err
		}

		info := ThreadFileInfo{Index: i}
		if looksLikeFileNode(node) {
			file, err := t.fileForPair(node)
			if err != nil {
				return nil, err
			}
			info.File = file

		} else {
			info.Links = make(Directory)
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
					info.Links[link.Name] = *file
				}
			}
		}

		files[i] = info
	}

	return files, nil
}

func (t *Textile) threadFile(block repo.Block) (*ThreadFilesInfo, error) {
	if block.Type != repo.FilesBlock {
		return nil, ErrBlockWrongType
	}

	files, err := t.fileAtTarget(block.Target)
	if err != nil {
		return nil, err
	}

	comments, err := t.ThreadComments(block.Id)
	if err != nil {
		return nil, err
	}

	likes, err := t.ThreadLikes(block.Id)
	if err != nil {
		return nil, err
	}

	threads := make([]string, 0)
	threads = t.fileThreads(block.Target)

	username, avatar := t.ContactDisplayInfo(block.AuthorId)

	return &ThreadFilesInfo{
		Block:    block.Id,
		Target:   block.Target,
		Date:     block.Date,
		AuthorId: block.AuthorId,
		Username: username,
		Avatar:   avatar,
		Caption:  block.Body,
		Files:    files,
		Comments: comments,
		Likes:    likes,
		Threads:  threads,
	}, nil
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
