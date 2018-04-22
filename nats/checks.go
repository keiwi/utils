package nats

import (
	"time"

	"github.com/keiwi/utils/models"
	"github.com/nats-io/go-nats"
	"gopkg.in/mgo.v2/bson"
)

func DeleteCheck(state *nats.Conn, data []byte) error {
	return state.Publish("checks.delete.send", data)
}

func FindCheck(state *nats.Conn, data []byte) ([]models.Check, error) {
	// msg, err := state.Request("checks.retrieve.find", data, time.Duration(viper.GetInt("nats.delay"))*time.Second)
	msg, err := state.Request("checks.retrieve.find", data, time.Duration(10)*time.Second)
	if err != nil {
		return nil, err
	}

	var checks []models.Check
	err = bson.UnmarshalJSON(msg.Data, &checks)
	if err != nil {
		return nil, err
	}
	return checks, nil
}

func HasCheck(state *nats.Conn, data []byte) (bool, error) {
	// msg, err := state.Request("checks.retrieve.has", data, time.Duration(viper.GetInt("nats.delay"))*time.Second)
	msg, err := state.Request("checks.retrieve.has", data, time.Duration(10)*time.Second)
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

func CreateCheck(state *nats.Conn, data []byte) error {
	return state.Publish("checks.create.send", data)
}

func UpdateCheck(state *nats.Conn, data []byte) error {
	return state.Publish("checks.update.send", data)
}
