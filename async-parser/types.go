package asyncparser

import (
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"sync"
)

type AppSettings struct {
	PathToSave   string
	DocumentName string
	IsForce      bool
	WorkersNum   int
	depth        int
}

type HandleParam struct {
	document *goquery.Document
	group    *sync.WaitGroup
	url      *url.URL
	links    chan DownloadParam
	depth    int
}

type DownloadParam struct {
	url            *url.URL
	remoteFileName string
	localFileName  string
	depth          int
}

type HandleDocument struct {
	filename string
	depth    int
}
