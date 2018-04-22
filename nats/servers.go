package nats

import (
	"time"

	"github.com/keiwi/utils/models"
	"github.com/nats-io/go-nats"
	"gopkg.in/mgo.v2/bson"
)

func DeleteServer(state *nats.Conn, data []byte) error {
	return state.Publish("servers.delete.send", data)
}

func FindServer(state *nats.Conn, data []byte) ([]models.Server, error) {
	// msg, err := state.Request("servers.retrieve.find", data, time.Duration(viper.GetInt("nats.delay"))*time.Second)
	msg, err := state.Request("servers.retrieve.find", data, time.Duration(10)*time.Second)
	if err != nil {
		return nil, err
	}

	var servers []models.Server
	err = bson.UnmarshalJSON(msg.Data, &servers)
	if err != nil {
		return nil, err
	}
	return servers, nil
}

func HasServer(state *nats.Conn, data []byte) (bool, error) {
	// msg, err := state.Request("servers.retrieve.has", data, time.Duration(viper.GetInt("nats.delay"))*time.Second)
	msg, err := state.Request("servers.retrieve.has", data, time.Duration(10)*time.Second)
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

func CreateServer(state *nats.Conn, data []byte) error {
	return state.Publish("servers.create.send", data)
}

func UpdateServer(state *nats.Conn, data []byte) error {
	return state.Publish("servers.update.send", data)
}
