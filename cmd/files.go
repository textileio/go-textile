package cmd

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	ipfspath "github.com/ipfs/go-path"
	"github.com/mitchellh/go-homedir"
	"github.com/textileio/go-textile/core"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/schema"
)

var errNothingToAdd = fmt.Errorf("nothing to add")

// ------------------------------------
// > files

const batchSize = 10

type millOpts struct {
	val map[string]string
}

func newMillOpts(ext map[string]string) millOpts {
	c := make(map[string]string)
	for k, v := range ext {
		c[k] = v
	}
	return millOpts{val: c}
}

func (m millOpts) setPlaintext(v bool) {
	m.val["plaintext"] = strconv.FormatBool(v)
}

func (m millOpts) setUse(v string) {
	m.val["use"] = v
}

// ------------------------------------
// > file add

func FileAdd(path string, threadID string, caption string, group bool, verbose bool) error {
	var pth string
	fi, err := os.Stdin.Stat()
	if err != nil {
		return err
	}

	// if stdin was not provided, then fetch the path manually
	if (fi.Mode() & os.ModeCharDevice) != 0 {
		if path == "" {
			return fmt.Errorf("neither stdin nor the argument 'path' were provided, try --help")
		}

		// check if path references a cid
		ipth, err := ipfspath.ParsePath(path)
		if err == nil {
			pth = ipth.String()
		} else {
			pth, err = homedir.Expand(path)
			if err != nil {
				pth = path
			}

			fi, err = os.Stat(path)
			if err != nil {
				return err
			}
		}
	}

	// fetch schema
	var thrd pb.Thread
	if _, err := executeJsonPbCmd(http.MethodGet, "threads/"+threadID, params{}, &thrd); err != nil {
		return err
	}

	if thrd.SchemaNode == nil {
		return core.ErrThreadSchemaRequired
	}

	var pths []string
	var dirs []*pb.Directory
	var count int

	start := time.Now()

	if fi != nil && fi.IsDir() {
		// add each file inside the directory
		err := filepath.Walk(pth, func(pth string, fi os.FileInfo, err error) error {
			if fi.IsDir() || fi.Name() == ".DS_Store" {
				return nil
			}
			pths = append(pths, pth)
			return nil
		})
		if err != nil {
			return err
		}
		msg := fmt.Sprintf("Found %d file", len(pths))
		if len(pths) != 1 {
			msg += "s"
		}
		output(msg)

		batches := batchPaths(pths, batchSize)
		for i, batch := range batches {

			ready := make(chan *pb.Directory, batchSize)
			go millBatch(batch, thrd.SchemaNode, ready, verbose)

			var cerr error
		loop:
			for {
				select {
				case dir, ok := <-ready:
					if !ok {
						break loop
					}

					if group == false {
						files, err := add(
							[]*pb.Directory{dir},
							threadID,
							strings.TrimSpace(fmt.Sprintf("%s (%d)", caption, count+1)),
							verbose)
						if err != nil {
							cerr = err
							break loop
						}

						output(fmt.Sprintf("File %d data=%s block=%s", count+1, files.Data, files.Block))
					} else {
						dirs = append(dirs, dir)
					}

					count++
				}
			}
			if cerr != nil {
				return cerr
			}

			output(fmt.Sprintf("Milled batch %d/%d", i+1, len(batches)))
		}

		if group && len(dirs) > 0 {
			files, err := add(dirs, threadID, caption, verbose)
			if err != nil {
				return err
			}
			output(fmt.Sprintf("Group data=%s block=%s", files.Data, files.Block))
		}

		if count == 0 {
			return errNothingToAdd
		}

	} else {
		// add the file
		dir, err := mill(pth, thrd.SchemaNode, verbose)
		if err != nil {
			return err
		}

		_, err = add([]*pb.Directory{dir}, threadID, caption, true)
		if err != nil {
			return err
		}

		count++
	}

	if verbose {
		dur := time.Now().Sub(start)
		msg := fmt.Sprintf("Added %d file", count)
		if count != 1 {
			msg += "s"
		}
		output(fmt.Sprintf("%s in %s", msg, dur.String()))
	}

	return nil
}

func add(dirs []*pb.Directory, threadID string, caption string, verbose bool) (*pb.Files, error) {
	data, err := pbMarshaler.MarshalToString(&pb.DirectoryList{Items: dirs})
	if err != nil {
		return nil, err
	}

	files := new(pb.Files)
	res, err := executeJsonPbCmd(http.MethodPost, "threads/"+threadID+"/files", params{
		opts:    map[string]string{"caption": caption},
		payload: strings.NewReader(data),
		ctype:   "application/json",
	}, files)
	if err != nil {
		return nil, err
	}

	if verbose {
		output(res)
	}

	return files, nil
}

func mill(pth string, node *pb.Node, verbose bool) (*pb.Directory, error) {
	ref, err := ipfspath.ParsePath(pth)
	if err == nil {
		pth = ref.String()
	}

	var f *os.File
	if ref == "" {
		if pth == "" {
			f = os.Stdin
		} else {
			var err error
			f, err = os.Open(pth)
			if err != nil {
				return nil, err
			}
			defer f.Close()

			fi, err := f.Stat()
			if err != nil {
				return nil, err
			}

			if fi.IsDir() {
				return nil, nil
			}
		}
	}

	var reader io.ReadSeeker
	var ctype string

	dir := &pb.Directory{Files: make(map[string]*pb.FileIndex)}

	// traverse the schema and collect generated files
	if node.Mill != "" {
		var res string
		file := &pb.FileIndex{}

		mopts := newMillOpts(node.Opts)
		mopts.setPlaintext(node.Plaintext)

		if node.Mill == "/json" {
			reader = f
			ctype = "application/json"
		} else if ref != "" {
			mopts.setUse(pth)
		} else {
			r, ct, err := multipartReader(f)
			if err != nil {
				return nil, err
			}
			reader = r
			ctype = ct
		}

		var err error
		res, file, err = handleStep(node.Mill, reader, mopts, ctype)
		if err != nil {
			return nil, err
		}

		if verbose {
			output(res)
		}

		dir.Files[schema.SingleFileTag] = file

	} else if len(node.Links) > 0 {

		// determine order
		steps, err := schema.Steps(node.Links)
		if err != nil {
			return nil, err
		}

		// send each link
		for _, step := range steps {
			var res string
			file := &pb.FileIndex{}

			mopts := newMillOpts(step.Link.Opts)
			mopts.setPlaintext(step.Link.Plaintext)

			if step.Link.Use == schema.FileTag {
				if reader != nil {
					_, _ = reader.Seek(0, 0)
				} else {
					if step.Link.Mill == "/json" {
						reader = f
						ctype = "application/json"
					} else if ref != "" {
						mopts.setUse(pth)
					} else {
						r, ct, err := multipartReader(f)
						if err != nil {
							return nil, err
						}
						reader = r
						ctype = ct
					}
				}

				res, file, err = handleStep(step.Link.Mill, reader, mopts, ctype)
				if err != nil {
					return nil, err
				}

			} else {
				if dir.Files[step.Link.Use].Hash == "" {
					return nil, fmt.Errorf(step.Link.Use + " not found")
				}
				mopts.setUse(dir.Files[step.Link.Use].Hash)

				res, err = executeJsonPbCmd(http.MethodPost, "mills"+step.Link.Mill, params{
					opts: mopts.val,
				}, file)
				if err != nil {
					return nil, err
				}
			}

			if verbose {
				output(res)
			}

			dir.Files[step.Name] = file
		}
	} else {
		return nil, schema.ErrEmptySchema
	}

	return dir, nil
}

func millBatch(pths []string, node *pb.Node, ready chan *pb.Directory, verbose bool) {
	wg := sync.WaitGroup{}

	for _, pth := range pths {
		wg.Add(1)

		go func(p string) {
			dir, err := mill(p, node, verbose)
			if err != nil {
				output("mill error: " + err.Error())
			} else {
				ready <- dir
			}
			wg.Done()
		}(pth)

	}

	wg.Wait()
	close(ready)
}

func batchPaths(pths []string, size int) [][]string {
	var batches [][]string

	for i := 0; i < len(pths); i += size {
		end := i + size

		if end > len(pths) {
			end = len(pths)
		}

		batches = append(batches, pths[i:end])
	}

	return batches
}

func handleStep(mil string, reader io.Reader, opts millOpts, ctype string) (string, *pb.FileIndex, error) {
	var file pb.FileIndex

	res, err := executeJsonPbCmd(http.MethodPost, "mills"+mil, params{
		opts:    opts.val,
		payload: reader,
		ctype:   ctype,
	}, &file)
	if err != nil {
		return "", nil, err
	}

	return res, &file, nil
}

func multipartReader(f *os.File) (io.ReadSeeker, string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filepath.Base(f.Name()))
	if err != nil {
		return nil, "", err
	}
	if _, err = io.Copy(part, f); err != nil {
		return nil, "", err
	}
	_ = writer.Close()
	return bytes.NewReader(body.Bytes()), writer.FormDataContentType(), nil
}

// ------------------------------------
// > file list thread

func FileListThread(threadID string, offset string, limit int) error {
	var list pb.FilesList
	res, err := executeJsonPbCmd(http.MethodGet, "files", params{opts: map[string]string{
		"thread": threadID,
		"offset": offset,
		"limit":  strconv.Itoa(limit),
	}}, &list)
	if err != nil {
		return err
	}
	if len(list.Items) > 0 {
		output(res)
	}

	if len(list.Items) < limit {
		return nil
	}

	if err := nextPage(); err != nil {
		return err
	}

	return FileListThread(threadID, list.Items[len(list.Items)-1].Block, limit)
}

// ------------------------------------
// > file list block

func FileListBlock(blockID string) error {
	// block
	urlPath := "blocks/" + blockID + "/files"

	// fetch
	res, err := executeJsonCmd(http.MethodGet, urlPath, params{}, nil)
	if err != nil {
		return err
	}
	output(res)

	// return
	return nil
}

// ------------------------------------
// > file get block

func FileGetBlock(blockID string, index int, path string, content bool) error {
	// block
	urlPath := "blocks/" + blockID + "/files"

	// index
	urlPath += "/" + strconv.Itoa(index)

	// path
	if path == "" {
		urlPath += "/."
	} else {
		urlPath += "/" + strings.Trim(path, "/")
	}

	// content
	if content {
		urlPath += "/content"
	} else {
		urlPath += "/meta"
	}

	// fetch
	if content {
		err := executeBlobCmd(http.MethodGet, urlPath, params{})
		if err != nil {
			return err
		}
	} else {
		res, err := executeJsonCmd(http.MethodGet, urlPath, params{}, nil)
		if err != nil {
			return err
		}
		output(res)
	}

	// return
	return nil
}

// ------------------------------------
// > file get

func FileGet(fileHash string, content bool) error {
	if content {
		err := executeBlobCmd(http.MethodGet, "file/"+fileHash+"/content", params{})
		return err
	} else {
		res, err := executeJsonCmd(http.MethodGet, "file/"+fileHash+"/meta", params{}, nil)
		if err != nil {
			return err
		}
		output(res)
	}
	return nil
}

// ------------------------------------
// > file ignore

func FileIgnore(blockID string) error {
	return BlockIgnore(blockID)
}

// ------------------------------------
// > files key

func FileKeys(dataID string) error {
	res, err := executeJsonCmd(http.MethodGet, "keys/"+dataID, params{}, nil)
	if err != nil {
		return err
	}
	output(res)
	return nil
}
