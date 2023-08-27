package asyncparser

import (
	"net/url"
)

const (
	PATH_TO_SAVE  = "./"
	DOCUMENT_NAME = "index.html"
	IS_FORCE      = true
	CHAN_SIZE     = 50
	WORKERS_NUM   = 5
)

func Start(rawUrl string, settings AppSettings) error {
	exactUrl, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}

	startPageParser := NewStartPageParser(exactUrl, settings)
	filepath, err := startPageParser.ParseStartPage()
	if err != nil {
		return err
	}

	linksChan := make(chan DownloadParam, CHAN_SIZE)

	worker := NewWorker(
		settings,
		exactUrl,
		linksChan,
	)

	for i := 0; i < settings.WorkersNum; i++ {
		go worker.Start()
	}

	linkHandler := NewLinkHandler(
		getLinksHandlers(),
		filepath,
		exactUrl,
		linksChan,
	)
	err = linkHandler.HandleDoc()
	if err != nil {
		return err
	}

	return nil
}

func getLinksHandlers() []func(param *HandleParam) {
	return []func(param *HandleParam){
		HandleHtml,
		HandleCss,
		HandleJs,
	}
}
