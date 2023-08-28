package asyncparser

import (
	"net/url"
)

const (
	PATH_TO_SAVE  = "./"
	DOCUMENT_NAME = "index.html"
	IS_FORCE      = false
	CHAN_SIZE     = 50
	WORKERS_NUM   = 5
)

func Start(rawUrl string, settings AppSettings) error {
	exactUrl, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}

	linksChan := make(chan DownloadParam, CHAN_SIZE)
	workersSemafor := make(chan bool)

	worker := NewWorker(
		settings,
		exactUrl,
		linksChan,
		NewLinkHandler(
			getLinksHandlers(),
			exactUrl,
			linksChan,
		),
		workersSemafor,
	)

	for i := 0; i < settings.WorkersNum; i++ {
		go worker.Start()
	}

	linksChan <- DownloadParam{
		exactUrl,
		"",
		"index.html",
		5,
	}

	select {
	case <-workersSemafor:
		break
	}
	close(linksChan)

	return nil
}

func getLinksHandlers() []func(param HandleParam) {
	return []func(param HandleParam){
		HandleHtml,
		HandleCss,
		HandleJs,
	}
}
