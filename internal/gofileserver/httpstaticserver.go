package gofileserver

import (
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pachirode/gofileserver/internal/pkg/log"
)

type HTTPStaticServer struct {
	Root     string
	Usage    string
	Upload   bool
	Delete   bool
	NoAccess bool
	Title    string
	Theme    string
	AuthType string
	Prefix   string
}

func (s *HTTPStaticServer) getRealPath(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	path = filepath.Clean(path)
	relativePath, err := filepath.Rel(s.Prefix, path)
	if err != nil {
		relativePath = path
	}
	realPath := filepath.Join(s.Root, relativePath)
	return filepath.ToSlash(realPath)
}

func (s *HTTPStaticServer) Index(ctx *gin.Context) {
	log.Infow("Index function called")
	path := ctx.Request.URL.Path
	realPath := s.getRealPath(path)
	println("============================")
	println(realPath)
	println("============================")
}
