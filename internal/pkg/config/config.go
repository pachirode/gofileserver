package config

import "os"

type Options struct {
	Addr            string
	AdminEmail      string
	AdminUsername   string
	AdminPassword   string
	Conf            *os.File
	Delete          bool
	Debug           bool
	GoogleTrackerID string
	HttpAuth        string
	NoIndex         bool
	NoAccess        bool
	Prefix          string
	Root            string
	SimpleAuth      bool
	Theme           string
	Title           string
	Upload          bool
	XHeaders        bool
}
