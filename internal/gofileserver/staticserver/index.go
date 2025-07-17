package staticserver

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/pachirode/gofileserver/internal/pkg/log"
)

type IndexFileItem struct {
	Path string
	Info os.FileInfo
}

func (s *HTTPStaticServer) makeIndex() error {
	indexes := make([]IndexFileItem, 0)
	err := filepath.Walk(s.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			println("Root: ", s.Root)
			log.Errorw("Visit path err", "err", err.Error())
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		path, _ = filepath.Rel(s.Root, path)
		path = filepath.ToSlash(path)
		indexes = append(indexes, IndexFileItem{Path: path, Info: info})
		return nil
	})

	s.indexes = indexes
	return err
}

func (s *HTTPStaticServer) Index(ctx *gin.Context) {
	log.Infow("Index function called")
	requestPath := ctx.Request.URL.Path
	realPath := s.getRealPath(requestPath)
	if ctx.Query("json") == "true" {
		s.JSONList(ctx)
		return
	}

	if ctx.Query("raw") == "false" || isDir(realPath) {
		renderHTML(ctx, "assets/index.html", s)
	}

	if ctx.Query("download") == "true" {
		s.Download(ctx)
	}
}

func (s *HTTPStaticServer) SysInfo(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"version": "unkonw",
	})
}
