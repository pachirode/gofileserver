package staticserver

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pachirode/gofileserver/internal/pkg/core"
	"github.com/pachirode/gofileserver/internal/pkg/log"
)

type Directory struct {
	sizeMap map[string]int64
	mutex   *sync.RWMutex
}

var dirSizeInfo = Directory{sizeMap: make(map[string]int64), mutex: &sync.RWMutex{}}

type HTTPFileInfo struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Type    string `json:"type"`
	Size    int64  `json:"size"`
	ModTime int64  `json:"mtime"`
}

func (s *HTTPStaticServer) JSONList(ctx *gin.Context) {
	log.Infow("Json list function called")
	requestPath := strings.Replace(ctx.Request.URL.Path, "/+", "", 1)
	realPath := s.getRealPath(requestPath)
	fileInfoMap := make(map[string]os.FileInfo, 0)

	search := ctx.Query("search")
	if search != "" {
		results := s.findIndex(search)
		if len(results) > 50 {
			results = results[:50]
		}
		for _, item := range results {
			if filepath.HasPrefix(item.Path, requestPath) {
				fileInfoMap[item.Path] = item.Info
			}
		}
	} else {
		infos, err := ioutil.ReadDir(realPath)
		if err != nil {
			core.WriteResponse(ctx, err, nil)
			return
		}
		for _, info := range infos {
			fileInfoMap[filepath.Join(requestPath, info.Name())] = info
		}
	}

	lrs := make([]HTTPFileInfo, 0)
	for path, info := range fileInfoMap {
		lr := HTTPFileInfo{
			Name:    info.Name(),
			Path:    path,
			ModTime: info.ModTime().UnixNano() / 1e6,
		}
		if search != "" {
			name, err := filepath.Rel(requestPath, path)
			if err != nil {
				log.Errorw("List File error", "err", err.Error())
			}
			lr.Name = filepath.ToSlash(name)
		}
		if info.IsDir() {
			name := deepPath(realPath, info.Name(), s.DeepPathMaxDepth)
			lr.Name = name
			lr.Path = filepath.Join(filepath.Dir(path), name)
			lr.Type = "dir"
			lr.Size = s.historyDirSizeInfo(lr.Path)
		} else {
			lr.Type = "file"
			lr.Size = info.Size()
		}
		lrs = append(lrs, lr)
	}

	ctx.JSON(200, gin.H{
		"files": lrs,
	})
}

func (s *HTTPStaticServer) getRealPath(path string) string {
	path = strings.Replace(path, "/+", "", 1)

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	path = filepath.Clean(path)
	relativePath, err := filepath.Rel(s.Prefix, path)
	if err != nil {
		relativePath = path
	}
	realPath := filepath.Join(s.Root, relativePath)
	return realPath
}

func (s *HTTPStaticServer) findIndex(text string) []IndexFileItem {
	ret := make([]IndexFileItem, 0)
	for _, item := range s.indexes {
		ok := true

		for _, keyword := range strings.Fields(text) {
			needContains := true
			if strings.HasPrefix(keyword, "-") {
				needContains = false
				keyword = keyword[1:]
			}

			if keyword == "" {
				continue
			}
			ok = (needContains == strings.Contains(strings.ToLower(item.Path), strings.ToLower(keyword)))
			if !ok {
				break
			}
		}
		if ok {
			ret = append(ret, item)
		}
	}
	return ret
}

func deepPath(basedir, name string, maxDepth int) string {
	for depth := 0; depth <= maxDepth; depth += 1 {
		finfos, err := ioutil.ReadDir(filepath.Join(basedir, name))
		if err != nil || len(finfos) != 1 {
			break
		}
		if finfos[0].IsDir() {
			name = filepath.ToSlash(filepath.Join(name, finfos[0].Name()))
		} else {
			break
		}
	}
	return name
}

func (s *HTTPStaticServer) historyDirSizeInfo(dir string) int64 {
	dirSizeInfo.mutex.Lock()
	size, ok := dirSizeInfo.sizeMap[dir]
	dirSizeInfo.mutex.Unlock()

	if ok {
		return size
	}

	for _, item := range s.indexes {
		if filepath.HasPrefix(item.Path, dir) {
			size += item.Info.Size()
		}
	}

	dirSizeInfo.mutex.Lock()
	dirSizeInfo.sizeMap[dir] = size
	dirSizeInfo.mutex.Unlock()

	return size
}
