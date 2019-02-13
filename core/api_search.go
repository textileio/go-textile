package core

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/textileio/textile-go/broadcast"
	"github.com/textileio/textile-go/pb"
)

func handleSearchStream(g *gin.Context, resultCh <-chan *pb.QueryResult, errCh <-chan error, cancel *broadcast.Broadcaster, events bool) {
	g.Stream(func(w io.Writer) bool {
		select {
		case <-g.Request.Context().Done():
			cancel.Close()

		case err := <-errCh:
			if events {
				g.SSEvent("error", err.Error())
			} else {
				g.String(http.StatusBadRequest, err.Error())
			}
			return false

		case res, ok := <-resultCh:
			if !ok {
				g.Status(http.StatusOK)
				return false
			}

			str, err := pbMarshaler.MarshalToString(res)
			if err != nil {
				g.String(http.StatusBadRequest, err.Error())
				break
			}

			if events {
				g.SSEvent("result", str)
			} else {
				g.Data(http.StatusOK, "application/json", []byte(str))
				g.Writer.Write([]byte("\n"))
			}
		}
		return true
	})
}
