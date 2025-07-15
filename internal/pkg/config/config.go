package config

import "os"

type Options struct {
	Conf          *os.File
	Addr          string
	Title         string
	AdminUsername string
	AdminPassword string
	AdminEmail    string
	Root          string
	HttpAuth      string
	SimpleAuth    bool
	Theme         string
	XHeaders      bool
	Upload        bool
	Delete        bool
	NoAccess      bool
	Debug         bool
}
