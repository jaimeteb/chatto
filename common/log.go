package common

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// SetLogger configures the logrus logger and sets the log level
func SetLogger() {
	if debug := os.Getenv("DEBUG"); debug == "true" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
}
