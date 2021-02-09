package logger

import (
	"os"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"

	log "github.com/sirupsen/logrus"
)

// SetLogger configures the logrus logger and sets the log level
func SetLogger() {
	if os.Getenv("DEBUG") == "true" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.SetFormatter(&runtime.Formatter{
		ChildFormatter: &log.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		},
		File: true,
		Line: true,
	})
}
