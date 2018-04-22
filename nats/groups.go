package nats

import (
	"time"

	"github.com/keiwi/utils/models"
	"github.com/nats-io/go-nats"
	"gopkg.in/mgo.v2/bson"
)

func DeleteGroup(state *nats.Conn, data []byte) error {
	return state.Publish("groups.delete.send", data)
}

func FindGroup(state *nats.Conn, data []byte) ([]models.Group, error) {
	// msg, err := state.Request("groups.retrieve.find", data, time.Duration(viper.GetInt("nats.delay"))*time.Second)
	msg, err := state.Request("groups.retrieve.find", data, time.Duration(10)*time.Second)
	if err != nil {
		return nil, err
	}

	var groups []models.Group
	err = bson.UnmarshalJSON(msg.Data, &groups)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func HasGroup(state *nats.Conn, data []byte) (bool, error) {
	// msg, err := state.Request("groups.retrieve.has", data, time.Duration(viper.GetInt("nats.delay"))*time.Second)
	msg, err := state.Request("groups.retrieve.has", data, time.Duration(10)*time.Second)
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

func CreateGroup(state *nats.Conn, data []byte) error {
	return state.Publish("groups.create.send", data)
}

func UpdateGroup(state *nats.Conn, data []byte) error {
	return state.Publish("groups.update.send", data)
}
