package asyncparser

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Worker struct {
	settings AppSettings
	url      *url.URL
	links    chan DownloadParam
}

func NewWorker(settings AppSettings, mainUrl *url.URL, links chan DownloadParam) *Worker {
	return &Worker{
		settings: settings,
		url:      mainUrl,
		links:    links,
	}
}

func (w Worker) Start() {
	for link := range w.links {
		func() {
			if w.isFileExist(link) && !w.settings.IsForce {
				return
			}

			reader, err := w.getDocument(link)
			if err != nil {
				log.Printf("cannot get document %s%s: %s", link.url, link.remoteFileName, err.Error())
				return
			}
			defer reader.Close()

			if err = w.saveDocument(reader, link); err != nil {
				log.Printf("cannot save document %s%s: %s", link.url, link.localFileName, err.Error())
				return
			}
		}()
	}
}

func (w Worker) isFileExist(doc DownloadParam) bool {
	path := filepath.Join(
		w.settings.PathToSave,
		w.url.Host,
		doc.url.Path,
		doc.localFileName,
	)
	_, err := os.Stat(path)
	return err == nil
}

func (w Worker) getDocument(d DownloadParam) (io.ReadCloser, error) {
	requestUrl := strings.Join([]string{d.url.String(), d.remoteFileName}, "/")
	response, err := http.Get(requestUrl)
	if err != nil {
		return nil, fmt.Errorf("http request error: %s", err.Error())
	}

	return response.Body, nil
}

func (w Worker) saveDocument(r io.ReadCloser, d DownloadParam) error {
	path := filepath.Join(
		w.settings.PathToSave,
		w.url.Host,
		d.url.Path,
	)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("cannot create directories: %s", err.Error())
	}

	documentPath := filepath.Join(
		path,
		d.localFileName,
	)
	fd, err := os.Create(documentPath)
	if err != nil {
		return fmt.Errorf("cannot create file: %s", err.Error())
	}
	defer fd.Close()

	if _, err := io.Copy(fd, r); err != nil {
		return fmt.Errorf("cannot write file: %s", err.Error())
	}

	return nil
}
