package staticserver

import (
	"time"

	"github.com/pachirode/gofileserver/internal/pkg/config"
	"github.com/pachirode/gofileserver/internal/pkg/log"
)

type HTTPStaticServer struct {
	Root             string
	Prefix           string
	Upload           bool
	Delete           bool
	Title            string
	Theme            string
	PlistProxy       string
	GoogleTrackerID  string
	AuthType         string
	DeepPathMaxDepth int
	NoIndex          bool

	indexes []IndexFileItem
}

func NewHTTPStaticServer(gcfg *config.Options) *HTTPStaticServer {
	s := &HTTPStaticServer{
		Root:            gcfg.Root,
		Theme:           gcfg.Theme,
		Title:           gcfg.Title,
		NoIndex:         gcfg.NoIndex,
		GoogleTrackerID: gcfg.GoogleTrackerID,
	}

	if !s.NoIndex {
		go func() {
			time.Sleep(1 * time.Second)
			for {
				log.Infow("Started marking search index")
				s.makeIndex()
				log.Infow("Completed search index")
				time.Sleep(time.Minute * 10)
			}
		}()
	}

	return s
}
