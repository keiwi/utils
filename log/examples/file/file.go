package main

import (
	"time"

	"github.com/keiwi/utils/log"
	"github.com/keiwi/utils/log/handlers/file"
)

func main() {
	l := log.NewLogger(log.DEBUG, []log.Reporter{file.NewFile(".", "%date%.log")})

	l.Debug("Testing debug")
	l.Info("Testing info")
	l.Warn("Testing warn")
	l.Error("Testing error")

	l.WithFields(log.Fields{"One field": "One value", "Second field": "Second value", "Hello": "World"}).Info("Multi line fields")

	l.WithField("test_field", "test_value").Debug("Testing with field debug")
	l.WithField("test_field", "test_value").Info("Testing with field info")
	l.WithField("test_field", "test_value").Warn("Testing with field warn")
	l.WithField("test_field", "test_value").Error("Testing with field error")

	time.Sleep(time.Second * 2)
}
