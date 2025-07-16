package staticserver

import (
	"embed"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pachirode/gofileserver/internal/pkg/log"
)

//go:embed assets
var assetsFS embed.FS

var (
	Assets = http.FS(assetsFS)
	_tmpls = make(map[string]*template.Template, 0)
)

func assetsFuncMap() template.FuncMap {
	return template.FuncMap{
		"title": strings.Title,
		"urlhash": func(path string) string {
			httpFile, err := Assets.Open(path)
			if err != nil {
				return path + "#no-such-file"
			}
			info, err := httpFile.Stat()
			if err != nil {
				return path + "#stat-error"
			}
			return fmt.Sprintf("%s?t=%d", path, info.ModTime().Unix())
		},
	}
}

func assetsContent(name string) string {
	fd, err := Assets.Open(name)
	if err != nil {
		log.Errorw("Assets Open template failed", "err", err.Error())
	}
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		log.Errorw("Assets read template failed", "err", err.Error())
	}
	return string(data)
}

func renderHTML(ctx *gin.Context, name string, v interface{}) {
	if tmpl, ok := _tmpls[name]; ok {
		if err := tmpl.Execute(ctx.Writer, v); err != nil {
			log.Errorw("Render HTML failed", "err", err.Error())
		}
		//	ctx.HTML(http.StatusOK, name, v)
		return
	}

	tmpl := template.Must(template.New(name).Funcs(assetsFuncMap()).Delims("[[", "]]").Parse(assetsContent(name)))
	_tmpls[name] = tmpl
	if err := tmpl.Execute(ctx.Writer, v); err != nil {
		log.Errorw("Render HTML failed", "err", err.Error())
	}
	// ctx.HTML(http.StatusOK, name, v)
}
