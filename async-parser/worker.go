package asyncparser

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Worker struct {
	settings    AppSettings
	url         *url.URL
	links       chan DownloadParam
	linkHandler *LinkHandler
	semafor     chan bool
}

func NewWorker(settings AppSettings, mainUrl *url.URL, links chan DownloadParam, linkHandler *LinkHandler, semafor chan bool) *Worker {
	return &Worker{
		settings:    settings,
		url:         mainUrl,
		links:       links,
		linkHandler: linkHandler,
		semafor:     semafor,
	}
}

func (w Worker) Start() {
	for link := range w.links {
		func() {
			slog.Info("download has started", "name", link.url.String()+link.remoteFileName)
			if w.isFileExist(link) && !w.settings.IsForce {
				slog.Warn("file already exist, use -force to ignore it", "file", link.url.String()+link.localFileName)
				return
			}

			reader, err := w.getDocument(link)
			if err != nil {
				slog.Error("cannot get document", "name", link.url.String()+link.remoteFileName, "err", err.Error())
				return
			}
			defer reader.Close()

			if err = w.saveDocument(reader, link); err != nil {
				slog.Error("cannot save document", "name", link.url.String()+link.remoteFileName, "err", err.Error())
				return
			}
			slog.Info("download has finished", "name", link.url.String()+link.remoteFileName)
		}()
	}
	w.semafor <- true
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
	d.url.RawQuery = ""
	d.url.Fragment = ""
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

	go w.linkHandler.HandleDoc(documentPath, d.depth)

	return nil
}
