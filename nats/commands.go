package nats

import (
	"time"

	"github.com/keiwi/utils/models"
	"github.com/nats-io/go-nats"
	"gopkg.in/mgo.v2/bson"
)

func DeleteCommand(state *nats.Conn, data []byte) error {
	return state.Publish("commands.delete.send", data)
}

func FindCommand(state *nats.Conn, data []byte) ([]models.Command, error) {
	// msg, err := state.Request("commands.retrieve.find", data, time.Duration(viper.GetInt("nats.delay"))*time.Second)
	msg, err := state.Request("commands.retrieve.find", data, time.Duration(10)*time.Second)
	if err != nil {
		return nil, err
	}

	var commands []models.Command
	err = bson.UnmarshalJSON(msg.Data, &commands)
	if err != nil {
		return nil, err
	}
	return commands, nil
}

func HasCommand(state *nats.Conn, data []byte) (bool, error) {
	// msg, err := state.Request("commands.retrieve.has", data, time.Duration(viper.GetInt("nats.delay"))*time.Second)
	msg, err := state.Request("commands.retrieve.has", data, time.Duration(10)*time.Second)
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

func CreateCommand(state *nats.Conn, data []byte) error {
	return state.Publish("commands.create.send", data)
}

func UpdateCommand(state *nats.Conn, data []byte) error {
	return state.Publish("commands.update.send", data)
}
