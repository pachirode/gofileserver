package staticserver

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/pachirode/gofileserver/internal/pkg/core"
	"github.com/pachirode/gofileserver/internal/pkg/log"
)

func (s *HTTPStaticServer) UploadOrMkdir(ctx *gin.Context) {
	log.Infow("Upload or make dir function called")
	requestPath := strings.Replace(ctx.Request.URL.Path, "/+", "", 1)
	dirPath := s.getRealPath(requestPath)

	file, err := ctx.FormFile("file")

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			core.WriteResponse(ctx, err, nil)
			return
		}
	}

	if file == nil {
		ctx.JSON(200, gin.H{
			"success":     true,
			"destination": dirPath,
		})
		return
	}

	fileName := file.Filename
	err = checkFilename(fileName)

	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	disPath := filepath.Join(dirPath, fileName)
	if err = ctx.SaveUploadedFile(file, disPath); err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	core.WriteResponse(ctx, nil, nil)
}
