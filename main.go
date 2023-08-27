package main

import (
	"flag"
	"fmt"
	asyncparser "go-parse/async-parser"
	"os"
)

func main() {
	pathToSave := flag.String("path", asyncparser.PATH_TO_SAVE, "path where program will store downloaded data")
	documentName := flag.String("docname", asyncparser.DOCUMENT_NAME, "name of the first downloaded file")
	isForce := flag.Bool("force", asyncparser.IS_FORCE, "force rewrite all documents")
	workersNum := flag.Int("workers", asyncparser.WORKERS_NUM, "number of workers used to parse")
	flag.Usage = func() {
		_, _ = fmt.Fprintln(os.Stderr, "Usage of ./go-parse [flags] url")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(os.Args) <= 1 {
		flag.Usage()
		os.Exit(1)
	}
	err := asyncparser.Start(os.Args[1], asyncparser.AppSettings{
		PathToSave:   *pathToSave,
		DocumentName: *documentName,
		IsForce:      *isForce,
		WorkersNum:   *workersNum,
	})
	if err != nil {
		println("Error while running application: ", err.Error())
		os.Exit(1)
	}
}
