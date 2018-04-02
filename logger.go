package main

import (
	"log"
	"os"
)

var (
	Trace *log.Logger
	Info  *log.Logger
)

func SetUpLoggers(
	traceHandler *os.File,
	infoHandler *os.File,
) {
	Trace = log.New(traceHandler,
		"TRACE : ",
		log.Ltime|log.Lshortfile)

	Info = log.New(infoHandler,
		"INFO  : ",
		log.Ltime|log.Lshortfile)

}
