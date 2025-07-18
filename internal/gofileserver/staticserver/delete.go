package staticserver

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pachirode/gofileserver/internal/pkg/core"
)

func (s *HTTPStaticServer) HDelete(ctx *gin.Context) {
	requestPath := strings.Replace(ctx.Request.URL.Path, "/+", "", 1)
	realPath := s.getRealPath(requestPath)

	err := os.RemoveAll(realPath)
	if err != nil {
		core.WriteResponse(ctx, err, nil)
	}
	core.WriteResponse(ctx, err, nil)
}
