package logger

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const timestampFormat = "2006-01-02 15:04:05.001 -0700 MST"

func init() {
	multiWriter := io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   "logs/server.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
	})
	log.SetOutput(multiWriter)
	log.SetLevel(log.DebugLevel)

	dateFormatter := &log.JSONFormatter{
		TimestampFormat: timestampFormat,
	}
	// output in JSON format
	log.SetFormatter(dateFormatter)
}
