package gofileserver

import (
	"github.com/gin-gonic/gin"

	"github.com/pachirode/gofileserver/internal/pkg/core"
	"github.com/pachirode/gofileserver/internal/pkg/errno"
	"github.com/pachirode/gofileserver/internal/pkg/log"
)

func installRouters(g *gin.Engine) error {
	g.NoRoute(func(ctx *gin.Context) {
		core.WriteResponse(ctx, errno.ErrPageNotFound, nil)
	})

	g.GET("/test", func(ctx *gin.Context) {
		log.C(ctx).Infow("Test function called")

		core.WriteResponse(ctx, nil, map[string]string{"status": "ok"})
	})

	return nil
}
