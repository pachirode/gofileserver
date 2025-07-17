package staticserver

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type Zip struct {
	*zip.Writer
}

func (z *Zip) Add(relPath, abspath string) error {
	info, rdc, err := z.statFile(abspath)
	if err != nil {
		return err
	}
	defer rdc.Close()

	hdr, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	hdr.Name = sanitizedName(relPath)
	if info.IsDir() {
		hdr.Name += "/"
	}
	hdr.Method = zip.Deflate
	writer, err := z.CreateHeader(hdr)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, rdc)
	return err
}

func (z *Zip) statFile(filename string) (info os.FileInfo, reader io.ReadCloser, err error) {
	info, err = os.Lstat(filename)
	if err != nil {
		return
	}

	if info.Mode()&os.ModeSymlink != 0 {
		var target string
		target, err = os.Readlink(filename)
		if err != nil {
			return
		}
		reader = ioutil.NopCloser(bytes.NewBuffer([]byte(target)))
	} else if !info.IsDir() {
		reader, err = os.Open(filename)
		if err != nil {
			return
		}
	} else {
		reader = ioutil.NopCloser(bytes.NewBuffer(nil))
	}
	return
}

func (s *HTTPStaticServer) Zip(ctx *gin.Context) {
	requestPath := strings.Replace(ctx.Request.URL.Path, "/+", "", 1)
	realDir := s.getRealPath(requestPath)
	zipFileName := filepath.Base(realDir) + ".zip"

	ctx.Header("Context-Type", "application/zip")
	ctx.Header("Context-Disposition", `"attachment; filename="`+zipFileName+`"`)

	zw := &Zip{Writer: zip.NewWriter(ctx.Writer)}
	defer zw.Close()

	filepath.Walk(realDir, func(path string, info os.FileInfo, err error) error {
		zipPath := path[len(realDir):]
		return zw.Add(zipPath, path)
	})
}
