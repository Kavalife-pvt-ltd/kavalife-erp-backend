package utils

import (
	"os"
	"strings"
	"sync"

	"github.com/paaart/kavalife-erp-backend/config"
	"github.com/sirupsen/logrus"
)

var (
	Log  *logrus.Logger
	once sync.Once
)

func InitLogger() *logrus.Logger {
	once.Do(func() {
		Log = logrus.New()
		appConfig := config.ConfigLoader()
		// Output
		Log.SetOutput(os.Stdout)

		// JSON formatter (best for prod + ELK + CloudWatch)
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})

		// Environment-based log level
		env := strings.ToLower(appConfig.APP_ENV)

		switch env {
		case "prod":
			Log.SetLevel(logrus.InfoLevel)
		default: // local, dev, test
			Log.SetLevel(logrus.DebugLevel)
		}

		Log.WithFields(logrus.Fields{
			"env": env,
		}).Info("Logger initialized")
	})

	return Log
}
