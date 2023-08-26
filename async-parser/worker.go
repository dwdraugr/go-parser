package asyncparser

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Worker struct {
	Settings AppSettings
	Url      *url.URL
	IsForce  bool
}

func (w Worker) Start(links <-chan DownloadParam) {
	for link := range links {
		if w.isFileExist(link) && !w.IsForce {
			continue
		}

		reader, err := w.getDocument(link)
		if err != nil {
			continue
		}
		defer reader.Close()

		if err = w.saveDocument(reader, link); err != nil {
			continue
		}
	}
}

func (w Worker) isFileExist(doc DownloadParam) bool {
	path := filepath.Join(
		w.Settings.PathToSave,
		w.Url.Host,
		doc.url.Path,
		doc.localFileName,
	)
	_, err := os.Stat(path)
	return err == nil
}

func (w Worker) getDocument(d DownloadParam) (io.ReadCloser, error) {
	requestUrl := strings.Join([]string{d.url.String(), d.remoteFileName}, "/")
	response, err := http.Get(requestUrl)
	fmt.Println(response.StatusCode)
	if err != nil {
		return nil, fmt.Errorf("http request error: %s", err.Error())
	}

	return response.Body, nil
}

func (w Worker) saveDocument(r io.ReadCloser, d DownloadParam) error {
	path := filepath.Join(
		w.Settings.PathToSave,
		w.Url.Host,
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

	println("saved ", documentPath)
	return nil
}
