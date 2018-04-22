package nats

import (
	"time"

	"github.com/keiwi/utils/models"
	"github.com/nats-io/go-nats"
	"gopkg.in/mgo.v2/bson"
)

func DeleteAlert(state *nats.Conn, data []byte) error {
	return state.Publish("alerts.delete.send", data)
}

func FindAlert(state *nats.Conn, data []byte) ([]models.Alert, error) {
	// msg, err := state.Request("alerts.retrieve.find", data, time.Duration(viper.GetInt("nats.delay"))*time.Second)
	msg, err := state.Request("alerts.retrieve.find", data, time.Duration(10)*time.Second)
	if err != nil {
		return nil, err
	}

	var alerts []models.Alert
	err = bson.UnmarshalJSON(msg.Data, &alerts)
	if err != nil {
		return nil, err
	}
	return alerts, nil
}

func HasAlert(state *nats.Conn, data []byte) (bool, error) {
	// msg, err := state.Request("alerts.retrieve.has", data, time.Duration(viper.GetInt("nats.delay"))*time.Second)
	msg, err := state.Request("alerts.retrieve.has", data, time.Duration(10)*time.Second)
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

func CreateAlert(state *nats.Conn, data []byte) error {
	return state.Publish("alerts.create.send", data)
}

func UpdateAlert(state *nats.Conn, data []byte) error {
	return state.Publish("alerts.update.send", data)
}
