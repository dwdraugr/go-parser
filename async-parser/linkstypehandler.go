package asyncparser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type LinkHandlerInterface interface {
	HandleDoc(string) error
}

type LinkHandler struct {
	handlers []func(*HandleParam)
	url      *url.URL
	links    chan<- DownloadParam
	wg       *sync.WaitGroup
}

func NewLinkHandler(
	handleFuncs []func(*HandleParam),
	url *url.URL,
	links chan<- DownloadParam,
	wg *sync.WaitGroup,
) LinkHandler {
	return LinkHandler{
		handlers: handleFuncs,
		url:      url,
		links:    links,
		wg:       wg,
	}
}

func (l LinkHandler) HandleDoc(filename string) error {
	_, err := os.Stat(filename)
	if filepath.Ext(filename) != PARSING_FILE_EXT || err == nil {
		return nil
	}

	fd, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("cannot open document: %s", err.Error())
	}
	defer fd.Close()

	doc, err := goquery.NewDocumentFromReader(fd)
	if err != nil {
		return fmt.Errorf("cannot handle file as html doc: %s", err.Error())
	}

	for _, handler := range l.handlers {
		l.wg.Add(1)
		go handler(&HandleParam{
			document: doc,
			group:    l.wg,
			url:      l.url,
			links:    l.links,
		})
	}

	return nil
}

func HandleHtml(p *HandleParam) {
	p.document.Find("a[href]").Each(func(i int, selection *goquery.Selection) {
		attr, _ := selection.Attr("href")
		if !checkIsInternalDomain(attr, p.url.Host) {
			return
		}
		p.links <- DownloadParam{
			generateFullUrl(attr, p.url),
			"",
			"index.html",
		}
	})
	p.group.Done()
}

func HandleCss(p *HandleParam) {
	p.document.Find("link[rel=stylesheet]").Each(func(i int, selection *goquery.Selection) {
		attr, _ := selection.Attr("href")
		if !checkIsInternalDomain(attr, p.url.Host) || attr == "" {
			return
		}
		uriElements := strings.Split(attr, "/")
		p.links <- DownloadParam{
			generateFullUrl(strings.Join(uriElements[:len(uriElements)-1], "/"), p.url),
			uriElements[len(uriElements)-1],
			uriElements[len(uriElements)-1],
		}
	})
	p.group.Done()
}

func HandleJs(p *HandleParam) {
	p.document.Find("script[src]").Each(func(i int, selection *goquery.Selection) {
		attr, _ := selection.Attr("src")
		if !checkIsInternalDomain(attr, p.url.Host) {
			return
		}
		uriElements := strings.Split(attr, "/")
		p.links <- DownloadParam{
			generateFullUrl(strings.Join(uriElements[:len(uriElements)-1], "/"), p.url),
			uriElements[len(uriElements)-1],
			uriElements[len(uriElements)-1],
		}
	})
	p.group.Done()
}

func checkIsInternalDomain(link string, domain string) bool {
	potentialUrl, err := url.Parse(link)
	if err != nil {
		return false
	}

	if potentialUrl.Host == domain || potentialUrl.Host == "" {
		return true
	}

	return false
}

func generateFullUrl(link string, mainUrl *url.URL) *url.URL {
	nextPageUrl, _ := url.Parse(link)
	nextPageUrl.Host = mainUrl.Host
	nextPageUrl.Scheme = mainUrl.Scheme

	return nextPageUrl
}
