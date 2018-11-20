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
	Caption  string              `json:"caption,omitempty"`
	Files    []ThreadFileInfo    `json:"files"`
	Comments []ThreadCommentInfo `json:"comments"`
	Likes    []ThreadLikeInfo    `json:"likes"`
	Threads  []string            `json:"threads"`
}

type ThreadCommentInfo struct {
	Id       string    `json:"id"`
	Date     time.Time `json:"date"`
	AuthorId string    `json:"author_id"`
	Username string    `json:"username,omitempty"`
	Body     string    `json:"body"`
}

type ThreadLikeInfo struct {
	Id       string    `json:"id"`
	Date     time.Time `json:"date"`
	AuthorId string    `json:"author_id"`
	Username string    `json:"username,omitempty"`
}

type ThreadJoinInfo struct {
	Id       string    `json:"id"`
	Date     time.Time `json:"date"`
	AuthorId string    `json:"author_id"`
	Username string    `json:"username,omitempty"`
}

type ThreadLeaveInfo struct {
	Id       string    `json:"id"`
	Date     time.Time `json:"date"`
	AuthorId string    `json:"author_id"`
	Username string    `json:"username,omitempty"`
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
	if len(blocks) == 0 {
		return list, nil
	}

	for _, block := range blocks {
		file, err := t.threadFile(&block)
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

	return t.threadFile(block)
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
		fnames := node.Links()

		i, err := strconv.Atoi(index.Name)
		if err != nil {
			return nil, err
		}

		info := ThreadFileInfo{Index: i}
		if len(fnames) > 0 {
			// directory of files
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
		} else {
			// single file
			file, err := t.fileForPair(node)
			if err != nil {
				return nil, err
			}
			info.File = file
		}

		files[i] = info
	}

	return files, nil
}

func (t *Textile) threadFile(block *repo.Block) (*ThreadFilesInfo, error) {
	files, err := t.fileAtTarget(block.Target)
	if err != nil {
		return nil, err
	}

	comments, err := t.fileComments(block.Id)
	if err != nil {
		return nil, err
	}

	likes, err := t.fileLikes(block.Id)
	if err != nil {
		return nil, err
	}

	threads := make([]string, 0)
	threads = t.fileThreads(block.Target)

	return &ThreadFilesInfo{
		Block:    block.Id,
		Target:   block.Target,
		Date:     block.Date,
		AuthorId: block.AuthorId,
		Username: t.ContactUsername(block.AuthorId),
		Caption:  block.Body,
		Files:    files,
		Comments: comments,
		Likes:    likes,
		Threads:  threads,
	}, nil
}

func (t *Textile) fileComments(target string) ([]ThreadCommentInfo, error) {
	comments := make([]ThreadCommentInfo, 0)

	query := fmt.Sprintf("type=%d and target='%s'", repo.CommentBlock, target)
	for _, block := range t.Blocks("", -1, query) {
		info := ThreadCommentInfo{
			Id:       block.Id,
			Date:     block.Date,
			AuthorId: block.AuthorId,
			Username: t.ContactUsername(block.AuthorId),
			Body:     block.Body,
		}
		comments = append(comments, info)
	}

	return comments, nil
}

func (t *Textile) fileLikes(target string) ([]ThreadLikeInfo, error) {
	likes := make([]ThreadLikeInfo, 0)

	query := fmt.Sprintf("type=%d and target='%s'", repo.LikeBlock, target)
	for _, block := range t.Blocks("", -1, query) {
		info := ThreadLikeInfo{
			Id:       block.Id,
			Date:     block.Date,
			AuthorId: block.AuthorId,
			Username: t.ContactUsername(block.AuthorId),
		}
		likes = append(likes, info)
	}

	return likes, nil
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
