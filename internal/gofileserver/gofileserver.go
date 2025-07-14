package gofileserver

import (
	"fmt"

	"github.com/marmotedu/Miniblog/pkg/version/verflag"
	"github.com/spf13/cobra"

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
	println("Start running")
	return nil
}
