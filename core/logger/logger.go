package logger

import (
	"errors"
	"log"
	"os"
	"strings"

	amitralog "github.com/amitrai48/logger"
)

type Interface interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type Logger struct {
	amitralog.Logger
}

type Config struct {
	File       string
	Level      string
	MuteStdout bool
}

var validLevel = map[string]bool{"debug": true, "info": true, "warn": true, "error": true, "fatal": true}

func New(conf Config) (Interface, error) {
	if conf.File == "" || !validLevel[strings.ToLower(conf.Level)] {
		return nil, errors.New("invalid logger config")
	}

	c := amitralog.Configuration{
		EnableConsole:     !conf.MuteStdout,
		ConsoleLevel:      amitralog.Fatal,
		ConsoleJSONFormat: false,
		EnableFile:        true,
		FileLevel:         strings.ToLower(conf.Level),
		FileJSONFormat:    true,
		FileLocation:      conf.File,
	}

	if err := amitralog.NewLogger(c, amitralog.InstanceZapLogger); err != nil {
		log.Fatalf("Could not instantiate log %s", err.Error())
	}
	l := amitralog.WithFields(amitralog.Fields{"hw": "15"})
	l.Infof("logger start successful")
	return l, nil
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Logger.Debugf(format, args)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, args)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Logger.Warnf(format, args)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, args)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Logger.Fatalf(format, args)
	os.Exit(2)
}
