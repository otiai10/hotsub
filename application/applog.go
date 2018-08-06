package application

import (
	"log"
	"regexp"
)

func applog(format string, v ...interface{}) {
	format = regexp.MustCompile("\n*$").ReplaceAllString(format, "\n")
	log.Printf("[COMMAND]\t"+format, v...)
}
