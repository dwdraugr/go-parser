package asyncparser

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Worker struct {
	settings   AppSettings
	links      <-chan DownloadParam
	wg         *sync.WaitGroup
	docHandler LinkHandlerInterface
}

func NewWorker(settings AppSettings, links chan DownloadParam, docHandler LinkHandlerInterface, wg *sync.WaitGroup) *Worker {
	return &Worker{
		settings:   settings,
		links:      links,
		docHandler: docHandler,
		wg:         wg,
	}
}

func (w Worker) Start() {
	for link := range w.links {
		func() {
			slog.Info("downloading has been started", "url", link.url.String()+link.remoteFileName)
			if w.isFileExist(link) && !w.settings.IsForce {
				slog.Info("file %s already downloaded", "url", link.url.String()+link.remoteFileName)
				return
			}
			reader, err := w.getDocument(link)
			if err != nil {
				slog.Error("cannot get document", "requestData", link, "err", err.Error())
				return
			}
			defer reader.Close()

			documentPath, err := w.saveDocument(reader, link)
			if err != nil {
				slog.Error("cannot save document", "requestData", link, "err", err.Error())
				return
			}
			slog.Info("downloading has finished", "url", link.url.String()+link.remoteFileName)
		}()
	}
	w.wg.Done()
}

func (w Worker) isFileExist(doc DownloadParam) bool {
	path := filepath.Join(
		w.settings.PathToSave,
		doc.url.Host,
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

func (w Worker) saveDocument(r io.ReadCloser, d DownloadParam) (string, error) {
	path := filepath.Join(
		w.settings.PathToSave,
		d.url.Host,
		d.url.Path,
	)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", fmt.Errorf("cannot create directories: %s", err.Error())
	}

	documentPath := filepath.Join(
		path,
		d.localFileName,
	)
	fd, err := os.Create(documentPath)
	if err != nil {
		return "", fmt.Errorf("cannot create file: %s", err.Error())
	}
	defer fd.Close()

	if _, err := io.Copy(fd, r); err != nil {
		return "", fmt.Errorf("cannot write file: %s", err.Error())
	}

	return documentPath, err
}
