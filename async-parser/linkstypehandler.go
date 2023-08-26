package asyncparser

import (
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"
)

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
