package staticserver

import (
	"github.com/gin-gonic/gin"

	"github.com/pachirode/gofileserver/internal/pkg/config"
	mw "github.com/pachirode/gofileserver/internal/pkg/middleware"
)

func InstallRouters(g *gin.Engine, gcfg *config.Options) error {
	mws := []gin.HandlerFunc{gin.Recovery(), mw.Cors, mw.NoCache}

	g.Use(mws...)

	ss := NewHTTPStaticServer(gcfg)

	g.GET("/", ss.Index)
	g.GET("/-/sysinfo", ss.SysInfo)
	g.GET("/-/assets/*filepath", ss.StaticFiles)
	g.GET("/+/*path", ss.Index)
	g.POST("/+/*path", ss.UploadOrMkdir)

	return nil
}
