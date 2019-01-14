package main

import (
	"log"
	"os"
)

var (
	trace *log.Logger
	info  *log.Logger
)

func setUpLoggers(
	traceHandler *os.File,
	infoHandler *os.File,
) {
	trace = log.New(traceHandler,
		"TRACE : ",
		log.Ltime|log.Lshortfile)

	info = log.New(infoHandler,
		"", log.Ltime)

}
