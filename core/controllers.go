package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/textileio/textile-go/crypto"
	"strings"
)

// gateway handles gateway http requests
func gateway(c *gin.Context) {
	path, contentType := parsePath(c.Param("path"))
	contentPath := fmt.Sprintf("%s/%s", c.Param("cid"), path)

	// look for block id
	blockId, exists := c.GetQuery("block")
	if exists {
		block, err := Node.Wallet.GetBlock(blockId)
		if err != nil {
			log.Errorf("error finding block %s: %s", blockId, err)
			return
		}
		thrd := Node.Wallet.GetThread(block.ThreadPubKey)
		if thrd == nil {
			log.Errorf("could not find thread for block: %s", block.Id)
			return
		}
		file, err := thrd.GetFileData(contentPath, block)
		if err != nil {
			log.Errorf("error decrypting path %s: %s", contentPath, err)
			c.Status(404)
			return
		}
		c.Data(200, contentType, file)
		return
	}

	// get raw file
	file, err := Node.Wallet.GetDataAtPath(contentPath)
	if err != nil {
		log.Errorf("error getting raw path %s: %s", contentPath, err)
		c.Status(404)
		return
	}

	// if key is provided, try to decrypt the file with it
	key, exists := c.GetQuery("key")
	if exists {
		plain, err := crypto.DecryptAES(file, []byte(key))
		if err != nil {
			log.Errorf("error decrypting %s: %s", contentPath, err)
			c.Status(404)
			return
		}
		c.Data(200, contentType, plain)
		return
	}

	// lastly, just return the raw bytes (standard gateway)
	c.Data(200, contentType, file)
}

// pin handles pin http requests
func pin(c *gin.Context) {
	c.Status(200)
}

func parsePath(path string) (parsed string, contentType string) {
	parts := strings.Split(path, ".")
	parsed = parts[0]
	if len(parts) == 1 {
		return parsed, "application/octet-stream"
	}
	switch parts[1] {
	case "jpg", "jpeg":
		return parsed, "image/jpeg"
	case "png":
		return parsed, "image/png"
	case "gif":
		return parsed, "image/gif"
	default:
		return parsed, "application/octet-stream"
	}
}
