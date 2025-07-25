package gofileserver

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/pachirode/gofileserver/internal/gofileserver/store"
	"github.com/pachirode/gofileserver/internal/pkg/config"
	"github.com/pachirode/gofileserver/internal/pkg/log"
	"github.com/pachirode/gofileserver/pkg/db"
)

var (
	recommandedHomeDir = ".config"
	defaultConfigName  = "gofileserver.yaml"
)

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(filepath.Join(home, recommandedHomeDir))
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(defaultConfigName)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Errorw("Failed to read viper configuration file", "err", err)
	}

	log.Debugw("Using configuration file", "file", viper.ConfigFileUsed())
}

func logOptions() *log.Options {
	return &log.Options{
		DisableCaller:     viper.GetBool("log.disable-caller"),
		DisableStacktrace: viper.GetBool("log.disable-stacktrace"),
		Level:             viper.GetString("log.level"),
		Format:            viper.GetString("log.format"),
		OutputPaths:       viper.GetStringSlice("log.output-paths"),
	}
}

func configurationOptions() *config.Options {
	return &config.Options{
		Addr:            viper.GetString("web.addr"),
		AdminUsername:   viper.GetString("web.admin_username"),
		AdminPassword:   viper.GetString("web.admin_password"),
		AdminEmail:      viper.GetString("web.admin_email"),
		Debug:           viper.GetBool("web.debug"),
		Delete:          viper.GetBool("web.delete"),
		GoogleTrackerID: viper.GetString("web.google_tracer_id"),
		HttpAuth:        viper.GetString("auth.http"),
		NoAccess:        viper.GetBool("web.noaccess"),
		NoIndex:         viper.GetBool("web.no_index"),
		Prefix:          viper.GetString("web.prefix"),
		Root:            viper.GetString("web.root"),
		SimpleAuth:      viper.GetBool("web.simpleauth"),
		Title:           viper.GetString("web.title"),
		Theme:           viper.GetString("web.theme"),
		Upload:          viper.GetBool("web.upload"),
		XHeaders:        viper.GetBool("web.xheaders"),
	}
}

func initStore() error {
	if !viper.GetBool("db.enable") {
		return nil
	}
	dbOptions := &db.MysqlOptions{
		Host:                  viper.GetString("db.host"),
		Username:              viper.GetString("db.username"),
		Password:              viper.GetString("db.password"),
		Database:              viper.GetString("db.database"),
		MaxIdleConnections:    viper.GetInt("db.max-idle-connections"),
		MaxOpenConnections:    viper.GetInt("db.max-open-connections"),
		MaxConnectionLifeTime: viper.GetDuration("db.max-connection-life-time"),
		LogLevel:              viper.GetInt("db.log-level"),
	}

	ins, err := db.NewMySQL(dbOptions)
	if err != nil {
		return err
	}

	_ = store.NewStore(ins)

	return nil
}
