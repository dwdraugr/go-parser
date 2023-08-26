package asyncparser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"os"
	"sync"
)

type LinkHandlerInterface interface {
}

type LinkHandler struct {
	Handlers []func(*HandleParam)
}

func (l LinkHandler) HandleDoc(filename string, pageUrl *url.URL, workersChan chan DownloadParam) error {
	var wg sync.WaitGroup
	fd, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("cannot open document: %s", err.Error())
	}
	defer fd.Close()

	doc, err := goquery.NewDocumentFromReader(fd)
	if err != nil {
		return fmt.Errorf("cannot handle file as html doc: %s", err.Error())
	}

	for _, handler := range l.Handlers {
		wg.Add(1)
		go handler(&HandleParam{
			document: doc,
			group:    &wg,
			url:      pageUrl,
			links:    workersChan,
		})
	}
	wg.Wait()

	return nil
}
