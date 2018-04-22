package nats

import (
	"time"

	"github.com/keiwi/utils/models"
	"github.com/nats-io/go-nats"
	"gopkg.in/mgo.v2/bson"
)

func DeleteUser(state *nats.Conn, data []byte) error {
	return state.Publish("users.delete.send", data)
}

func FindUser(state *nats.Conn, data []byte) ([]models.User, error) {
	// msg, err := state.Request("users.retrieve.find", data, time.Duration(viper.GetInt("nats.delay"))*time.Second)
	msg, err := state.Request("users.retrieve.find", data, time.Duration(10)*time.Second)
	if err != nil {
		return nil, err
	}

	var users []models.User
	err = bson.UnmarshalJSON(msg.Data, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func HasUser(state *nats.Conn, data []byte) (bool, error) {
	// msg, err := state.Request("users.retrieve.has", data, time.Duration(viper.GetInt("nats.delay"))*time.Second)
	msg, err := state.Request("users.retrieve.has", data, time.Duration(10)*time.Second)
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

func CreateUser(state *nats.Conn, data []byte) error {
	return state.Publish("users.create.send", data)
}

func UpdateUser(state *nats.Conn, data []byte) error {
	return state.Publish("users.update.send", data)
}
