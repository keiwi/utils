package nats

import (
	"time"

	"github.com/keiwi/utils/models"
	"github.com/nats-io/go-nats"
	"gopkg.in/mgo.v2/bson"
)

func DeleteAlertOption(state *nats.Conn, data []byte) error {
	return state.Publish("alert_options.delete.send", data)
}

func FindAlertOption(state *nats.Conn, data []byte) ([]models.AlertOption, error) {
	// msg, err := state.Request("alert_options.retrieve.find", data, time.Duration(viper.GetInt("nats.delay"))*time.Second)
	msg, err := state.Request("alert_options.retrieve.find", data, time.Duration(10)*time.Second)
	if err != nil {
		return nil, err
	}

	var alertOptions []models.AlertOption
	err = bson.UnmarshalJSON(msg.Data, &alertOptions)
	if err != nil {
		return nil, err
	}
	return alertOptions, nil
}

func HasAlertOption(state *nats.Conn, data []byte) (bool, error) {
	// msg, err := state.Request("alert_options.retrieve.has", data, time.Duration(viper.GetInt("nats.delay"))*time.Second)
	msg, err := state.Request("alert_options.retrieve.has", data, time.Duration(10)*time.Second)
	if err != nil {
		return false, err
	}

	var has bool
	err = bson.UnmarshalJSON(msg.Data, &has)
	if err != nil {
		return false, err
	}
	return has, nil
}

func CreateAlertOption(state *nats.Conn, data []byte) error {
	return state.Publish("alert_options.create.send", data)
}

func UpdateAlertOption(state *nats.Conn, data []byte) error {
	return state.Publish("alert_options.update.send", data)
}
