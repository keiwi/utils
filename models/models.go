package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Model struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}

// Alert struct
type Alert struct {
	Model     `bson:",inline"`
	AlertID   bson.ObjectId `json:"alert_id" bson:"alert_id"`
	ClientID  bson.ObjectId `json:"client_id" bson:"client_id"`
	Timestamp time.Time     `json:"timestamp" bson:"timestamp"`
	Value     string        `json:"value" bson:"value"`
}

// AlertOption struct
type AlertOption struct {
	Model     `bson:",inline"`
	ClientID  bson.ObjectId `json:"client_id" bson:"client_id"`
	CommandID bson.ObjectId `json:"command_id" bson:"command_id"`
	Alert     string        `json:"alert" bson:"alert"`
	Value     string        `json:"value" bson:"value"`
	Count     int           `json:"count" bson:"count"`
	Delay     int           `json:"delay" bson:"delay"`
	Service   string        `json:"service" bson:"service"`
}

// Check struct
type Check struct {
	Model     `bson:",inline"`
	CommandID bson.ObjectId `json:"command_id" bson:"command_id"`
	ClientID  bson.ObjectId `json:"client_id" bson:"client_id"`
	Response  string        `json:"response" bson:"response"`
	Checked   bool          `json:"checked" bson:"checked"`
	Error     bool          `json:"error" bson:"error"`
	Finished  bool          `json:"finished" bson:"finished"`
}

// Client struct
type Client struct {
	Model    `bson:",inline"`
	GroupIDs []bson.ObjectId `json:"group_ids" bson:"group_ids"`
	IP       string          `json:"ip" bson:"ip"`
	Name     string          `json:"name" bson:"name"`
}

// Command struct
type Command struct {
	Model       `bson:",inline"`
	Command     string `json:"command" bson:"command"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Format      string `json:"format" bson:"format"`
}

// Group struct
type Group struct {
	Model    `bson:",inline"`
	Commands []GroupCommand `json:"commands" bson:"commands"`
	Name     string         `json:"name" bson:"name"`
}

// GroupCommand struct
type GroupCommand struct {
	ID        bson.ObjectId `json:"id" bson:"id,omitempty"`
	CommandID bson.ObjectId `json:"command_id" bson:"command_id"`
	NextCheck int           `json:"next_check" bson:"next_check"`
	StopError bool          `json:"stop_error" bson:"stop_error"`
}

// Server struct
type Server struct {
	Model `bson:",inline"`
	IP    string `json:"ip" bson:"ip"`
	Name  string `json:"name" bson:"name"`
}

// Upload struct
// Upload struct
type Upload struct {
	Model         `bson:",inline" bson:"created_at"`
	Name          string `json:"name" bson:"name"`
	Checksum      string `json:"checksum" bson:"checksum"`
	Version       string `json:"version" bson:"version"`
	Patch         bool   `json:"patch" bson:"patch"`
	PatchChecksum string `json:"patch_checksum" bson:"patch_checksum"`
}

// User struct
type User struct {
	Model    `bson:",inline"`
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}
