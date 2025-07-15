package gofileserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/marmotedu/Miniblog/pkg/version/verflag"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/pachirode/gofileserver/internal/pkg/log"
)

var cfgFile string

func NewGofileserverCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "gofileserver",
		Short:        "A mini go file server",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Init(logOptions())
			defer log.Sync()

			return run()
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, go %q", cmd.CommandPath(), args)
				}
			}

			return nil
		},
	}

	cobra.OnInitialize(initConfig)

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "The path to the go file server configurate.")
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	verflag.AddFlags(cmd.PersistentFlags())

	return cmd
}

func run() error {
	if err := initStore(); err != nil {
		log.Errorw("Connect db failed", "err", err.Error())
		return nil
	}

	g := gin.New()

	if err := installRouters(g); err != nil {
		return err
	}

	httpServer := startInsecureServer(g)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Infow("Shutting down server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Errorw("Insecure server forced to shutdown", "err", err)
		return err
	}

	log.Infow("Server shutdown")

	return nil
}

func startInsecureServer(g *gin.Engine) *http.Server {
	httpServer := &http.Server{Addr: viper.GetString("web.addr"), Handler: g}

	log.Infow("Start to listening the incoming requests on http address", "addr", viper.GetString("web.addr"))
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatalw(err.Error())
		}
	}()

	return httpServer
}
