package main

import (
	"errors"

	"github.com/keiwi/utils/log"
	"github.com/keiwi/utils/log/handlers/cli"
)

func main() {
	l := log.NewLogger(log.DEBUG, []log.Reporter{cli.NewCli()})

	l.Debug("Testing debug")
	l.Info("Testing info")
	l.Warn("Testing warn")
	l.Error("Testing error")

	l.WithFields(log.Fields{"One field": "One value", "Second field": "Second value", "Hello": "World"}).Info("Multi line fields")

	l.WithField("test_field", "test_value").Debug("Testing with field debug")
	l.WithField("test_field", "test_value").Info("Testing with field info")
	l.WithField("test_field", "test_value").Warn("Testing with field warn")
	l.WithField("test_field", "test_value").Error("Testing with field error")

	l.WithError(errors.New("test error")).Error("error testing")
}
