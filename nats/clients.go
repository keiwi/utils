package nats

import (
	"time"

	"github.com/keiwi/utils/models"
	"github.com/nats-io/go-nats"
	"gopkg.in/mgo.v2/bson"
)

func DeleteClient(state *nats.Conn, data []byte) error {
	return state.Publish("clients.delete.send", data)
}

func FindClient(state *nats.Conn, data []byte) ([]models.Client, error) {
	// msg, err := state.Request("clients.retrieve.find", data, time.Duration(viper.GetInt("nats.delay"))*time.Second)
	msg, err := state.Request("clients.retrieve.find", data, time.Duration(10)*time.Second)
	if err != nil {
		return nil, err
	}

	var clients []models.Client
	err = bson.UnmarshalJSON(msg.Data, &clients)
	if err != nil {
		return nil, err
	}
	return clients, nil
}

func HasClient(state *nats.Conn, data []byte) (bool, error) {
	// msg, err := state.Request("clients.retrieve.has", data, time.Duration(viper.GetInt("nats.delay"))*time.Second)
	msg, err := state.Request("clients.retrieve.has", data, time.Duration(10)*time.Second)
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

func CreateClient(state *nats.Conn, data []byte) error {
	return state.Publish("clients.create.send", data)
}

func UpdateClient(state *nats.Conn, data []byte) error {
	return state.Publish("clients.update.send", data)
}
