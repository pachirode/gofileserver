package staticserver

import (
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pachirode/gofileserver/internal/pkg/log"
)

func (s *HTTPStaticServer) Download(ctx *gin.Context) {
	log.Infow("Download function called")
	requestPath := ctx.Request.URL.Path
	realPath := s.getRealPath(requestPath)
	ctx.Header("Content-Disposition", "attachment; filename="+strconv.Quote(filepath.Base(requestPath)))
	ctx.File(realPath)
}
