package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

// GetLogger returns a new logger
func GetLogger(debug bool) *logrus.Logger {
	log := &logrus.Logger{
		Out:   os.Stdout,
		Level: logrus.InfoLevel,
		Formatter: &logrus.TextFormatter{
			TimestampFormat:  "2006-01-02 15:04:05",
			FullTimestamp:    true,
			ForceColors:      true,
			QuoteEmptyFields: true,
		},
	}
	if debug {
		log.SetLevel(logrus.DebugLevel)
	}
	return log
}
