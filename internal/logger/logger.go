package logger

import (
	"os"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"

	log "github.com/sirupsen/logrus"
)

// SetLogger configures the logrus logger and sets the log level
func SetLogger(debug bool) {
	if debug || os.Getenv("DEBUG") == "true" {
		log.SetLevel(log.DebugLevel)
		log.SetFormatter(&runtime.Formatter{
			ChildFormatter: &log.TextFormatter{
				TimestampFormat: "2006-01-02 15:04:05.000000",
				FullTimestamp:   true,
			},
			File: true,
			Line: true,
		})
	} else {
		log.SetLevel(log.InfoLevel)
		log.SetFormatter(&log.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05.000000",
			FullTimestamp:   true,
		})
	}
}
