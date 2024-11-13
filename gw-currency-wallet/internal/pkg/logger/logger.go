package logger

import (
    "os"
    "github.com/sirupsen/logrus"
)

type LoggerConfig struct {
	Level logrus.Level
	LogToFile bool
	LogFilePath string
}

type Logger struct {
    logger *logrus.Logger
}

func NewLogger(config LoggerConfig) *Logger {
    log := logrus.New()
    log.SetLevel(config.Level)
    log.SetFormatter(&logrus.JSONFormatter{})

    if config.LogToFile {
        logFile, err := os.OpenFile(config.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
        if err == nil {
            log.SetOutput(logFile)
        } else {
            log.SetOutput(os.Stdout)
            log.Warn("Failed to open log file, logging to stdout instead")
        }
    } else {
        log.SetOutput(os.Stdout)
    }

    return &Logger{logger: log}
}

func (l *Logger) Info(msg string) {
    l.logger.Info(msg)
}

// ...