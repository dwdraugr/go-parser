package asyncparser

import (
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"sync"
)

type AppSettings struct {
	PathToSave        string
	DocumentName      string
	DocumentExtension string
}

type HandleParam struct {
	document *goquery.Document
	group    *sync.WaitGroup
	url      *url.URL
	links    chan DownloadParam
}

type DownloadParam struct {
	url            *url.URL
	remoteFileName string
	localFileName  string
}
