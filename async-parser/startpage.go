package asyncparser

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type StartPageParserInterface interface {
	ParseStartPage() (string, error)
}

type StartPageParser struct {
	targetUrl *url.URL
	settings  AppSettings
}

func NewStartPageParser(targetUrl *url.URL, settings AppSettings) *StartPageParser {
	return &StartPageParser{
		targetUrl: targetUrl,
		settings:  settings,
	}
}

func (s StartPageParser) ParseStartPage() (string, error) {
	reader, err := s.getPage()
	if err != nil {
		return "", fmt.Errorf("error while getting page: %s", err.Error())
	}
	defer reader.Close()

	filePath, err := s.savePage(reader)
	if err != nil {
		return "", fmt.Errorf("error while saving page: %s", err.Error())
	}

	return filePath, nil
}

func (s StartPageParser) getPage() (io.ReadCloser, error) {
	response, err := http.Get(s.targetUrl.String())
	if err != nil {
		return nil, fmt.Errorf("http request error: %s", err.Error())
	}

	return response.Body, nil
}

func (s StartPageParser) savePage(reader io.ReadCloser) (string, error) {
	path := filepath.Join(
		s.settings.PathToSave,
		s.targetUrl.Host,
		s.targetUrl.Path,
	)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", fmt.Errorf("cannot create directories: %s", err.Error())
	}

	documentPath := filepath.Join(
		path,
		s.settings.DocumentName,
	)
	fd, err := os.Create(documentPath)
	if err != nil {
		return "", fmt.Errorf("cannot create file: %s", err.Error())
	}
	defer fd.Close()

	if _, err := io.Copy(fd, reader); err != nil {
		return "", fmt.Errorf("cannot write file: %s", err.Error())
	}

	return documentPath, nil
}
