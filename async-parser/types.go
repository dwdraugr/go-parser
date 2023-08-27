package asyncparser

import (
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"sync"
)

const PARSING_FILE_EXT = ".html"

type AppSettings struct {
	PathToSave   string
	DocumentName string
	IsForce      bool
	WorkersNum   int
}

type HandleParam struct {
	document *goquery.Document
	group    *sync.WaitGroup
	url      *url.URL
	links    chan<- DownloadParam
}

type DownloadParam struct {
	url            *url.URL
	remoteFileName string
	localFileName  string
}
