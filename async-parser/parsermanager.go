package asyncparser

import (
	"net/url"
	"sync"
)

const (
	PATH_TO_SAVE  = "./"
	DOCUMENT_NAME = "index.html"
	IS_FORCE      = false
	CHAN_SIZE     = 50
	WORKERS_NUM   = 5
	DEPTH         = 5
)

func Start(rawUrl string, settings AppSettings) error {
	var wg sync.WaitGroup
	exactUrl, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}

	linksChan := make(chan DownloadParam, CHAN_SIZE)
	linkHandler := NewLinkHandler(
		getLinksHandlers(),
		exactUrl,
		linksChan,
		&wg,
	)
	worker := NewWorker(
		settings,
		linksChan,
		linkHandler,
		&wg,
	)

	for i := 0; i < settings.WorkersNum; i++ {
		go worker.Start()
	}

	linksChan <- DownloadParam{
		url:            exactUrl,
		remoteFileName: "",
		localFileName:  settings.DocumentName,
		depth:          DEPTH,
	}

	wg.Wait()
	close(linksChan)

	return nil
}

func getLinksHandlers() []func(param *HandleParam) {
	return []func(param *HandleParam){
		HandleHtml,
		HandleCss,
		HandleJs,
	}
}
