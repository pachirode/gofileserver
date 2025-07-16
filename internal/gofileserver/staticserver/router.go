package staticserver

import (
	"github.com/gin-gonic/gin"

	"github.com/pachirode/gofileserver/internal/pkg/config"
	"github.com/pachirode/gofileserver/internal/pkg/core"
	"github.com/pachirode/gofileserver/internal/pkg/log"
)

func InstallRouters(g *gin.Engine, gcfg *config.Options) error {
	g.StaticFS("/-/", Assets)

	g.GET("/test", func(ctx *gin.Context) {
		log.C(ctx).Infow("Test function called")

		core.WriteResponse(ctx, nil, map[string]string{"status": "ok"})
	})

	ss := NewHTTPStaticServer(gcfg)

	g.GET("/", ss.Index)

	return nil
}
