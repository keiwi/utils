package nats

import (
	"time"

	"github.com/keiwi/utils/models"
	"github.com/nats-io/go-nats"
	"gopkg.in/mgo.v2/bson"
)

func DeleteUpload(state *nats.Conn, data []byte) error {
	return state.Publish("uploads.delete.send", data)
}

func FindUpload(state *nats.Conn, data []byte) ([]models.Upload, error) {
	// msg, err := state.Request("uploads.retrieve.find", data, time.Duration(viper.GetInt("nats.delay"))*time.Second)
	msg, err := state.Request("uploads.retrieve.find", data, time.Duration(10)*time.Second)
	if err != nil {
		return nil, err
	}

	var uploads []models.Upload
	err = bson.UnmarshalJSON(msg.Data, &uploads)
	if err != nil {
		return nil, err
	}
	return uploads, nil
}

func HasUpload(state *nats.Conn, data []byte) (bool, error) {
	// msg, err := state.Request("uploads.retrieve.has", data, time.Duration(viper.GetInt("nats.delay"))*time.Second)
	msg, err := state.Request("uploads.retrieve.has", data, time.Duration(10)*time.Second)
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

func CreateUpload(state *nats.Conn, data []byte) error {
	return state.Publish("uploads.create.send", data)
}

func UpdateUpload(state *nats.Conn, data []byte) error {
	return state.Publish("uploads.update.send", data)
}
