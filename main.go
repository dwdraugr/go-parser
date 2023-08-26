package main

import (
	asyncparser "go-parse/async-parser"
	"net/url"
	"os"
)

const (
	ENV_PATH_TO_SAVE  = "GOPARSE_PATH_TO_SAVE"
	ENV_DOCUMENT_NAME = "GOPARSE_DOCUMENT_NAME"
)

func main() {
	if len(os.Args) <= 1 {
		println("There is no URL provided")
		os.Exit(1)
	}
	rawUrl := os.Args[1]
	exactUrl, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		println("Invalid URL", err)
		os.Exit(1)
	}

	settings := asyncparser.AppSettings{
		PathToSave:   getEnv(ENV_PATH_TO_SAVE, "./"),
		DocumentName: getEnv(ENV_DOCUMENT_NAME, "index.html"),
	}

	sp := asyncparser.NewStartPageParser(exactUrl, settings)

	filepath, err := sp.ParseStartPage()
	if err != nil {
		panic(err)
	}

	handlers := []func(*asyncparser.HandleParam){
		asyncparser.HandleHtml,
		asyncparser.HandleCss,
		asyncparser.HandleJs,
	}
	workersChan := make(chan asyncparser.DownloadParam)
	worker := &asyncparser.Worker{
		Settings: settings,
		Url:      exactUrl,
		IsForce:  false,
	}

	for i := 0; i < 1; i++ {
		go worker.Start(workersChan)
	}

	linkHandler := &asyncparser.LinkHandler{Handlers: handlers}
	err = linkHandler.HandleDoc(filepath, exactUrl, workersChan)
	if err != nil {
		panic(err)
	}
}

func getEnv(key string, defaultValue string) string {
	value, isOk := os.LookupEnv(key)
	if isOk {
		return value
	}

	return defaultValue
}
