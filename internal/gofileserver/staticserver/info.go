package staticserver

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/shogo82148/androidbinary/apk"

	"github.com/gin-gonic/gin"
	"github.com/pachirode/gofileserver/internal/pkg/core"
	"github.com/pachirode/gofileserver/internal/pkg/log"
)

type FileJSONInfo struct {
	Name    string      `json:"name"`
	Type    string      `json:"type"`
	Size    int64       `json:"size"`
	Path    string      `json:"path"`
	ModTime int64       `json:"mtime"`
	Extra   interface{} `json:"extra,omitempty"`
}

type ApkInfo struct {
	PackageName  string `json:"packageName"`
	MainActivity string `json:"mainActivity"`
	Version      struct {
		Code int    `json:"code"`
		Name string `json:"name"`
	} `json:"version"`
}

func (s *HTTPStaticServer) Info(ctx *gin.Context) {
	requestPath := strings.Replace(ctx.Request.URL.Path, "/+", "", 1)
	realPath := s.getRealPath(requestPath)

	fi, err := os.Stat(realPath)
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}
	fileJsonInfo := &FileJSONInfo{
		Name:    fi.Name(),
		Size:    fi.Size(),
		Path:    requestPath,
		ModTime: fi.ModTime().UnixNano() / 1e6,
	}
	ext := filepath.Ext(requestPath)
	switch ext {
	case ".md":
		fileJsonInfo.Type = "markdown"
	case ".apk":
		fileJsonInfo.Type = "apk"
		fileJsonInfo.Extra = parseApkInfo(realPath)
	case "":
		fileJsonInfo.Type = "dir"
	default:
		fileJsonInfo.Type = "text"
	}

	ctx.JSON(200, fileJsonInfo)
}

func parseApkInfo(path string) (apkInfo *ApkInfo) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorw("parse-apk-info panic", "err", err)
		}
	}()

	apkFile, err := apk.OpenFile(path)
	if err != nil {
		return
	}

	apkInfo = &ApkInfo{}
	apkInfo.MainActivity, _ = apkFile.MainActivity()
	apkInfo.PackageName = apkFile.PackageName()
	apkInfo.Version.Code = apkFile.Manifest().VersionCode
	apkInfo.Version.Name = apkFile.Manifest().VersionName
	return
}
