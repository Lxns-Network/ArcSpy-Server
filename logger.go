package main

import (
	"fmt"
	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"strings"
)

type CustomFormatter struct {
	log.TextFormatter
}

func (f *CustomFormatter) Format(entry *log.Entry) ([]byte, error) {
	var levelColor int
	switch entry.Level {
	case log.DebugLevel, log.TraceLevel:
		levelColor = 36 // cyan
	case log.WarnLevel:
		levelColor = 33 // yellow
	case log.ErrorLevel, log.FatalLevel, log.PanicLevel:
		levelColor = 31 // red
	default:
		levelColor = 0 // white
	}
	return []byte(fmt.Sprintf("\x1b[%dm[%s] [%s] ▶ \u001B[0m%s\n", levelColor, entry.Time.Format(f.TimestampFormat), strings.ToUpper(entry.Level.String()), entry.Message)), nil
}

func init() {
	log.SetFormatter(&CustomFormatter{log.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		ForceColors:            true,
		DisableLevelTruncation: true,
	}})
	log.SetOutput(colorable.NewColorableStdout())
	// log.SetLevel(log.DebugLevel)
}
