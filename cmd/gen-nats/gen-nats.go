package main

import (
	"flag"
	"html/template"
	"os"
	"path"
	"strings"
	"unicode"
)

var funcMap = map[string]interface{}{
	"CamelCase": func(s string) string {
		return toCamelCase(s, true)
	},
	"LowerCamelCase": func(s string) string {
		return toCamelCase(s, false)
	},
}

func toCamelCase(s string, upper bool) string {
	prev := 'a'
	if upper {
		prev = ' '
	}

	s = strings.Map(
		func(r rune) rune {
			if isSeparator(prev, true) {
				prev = r
				return unicode.ToTitle(r)
			}

			prev = r
			return r
		},
		s)

	return strings.Replace(s, "_", "", -1)
}

func isSeparator(r rune, underscore bool) bool {
	if r <= 0x7F {
		switch {
		case '0' <= r && r <= '9':
			return false
		case 'a' <= r && r <= 'z':
			return false
		case 'A' <= r && r <= 'Z':
			return false
		case r == '_':
			return underscore
		}
		return true
	}
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return false
	}
	return unicode.IsSpace(r)
}

var databaseTemplate = template.Must(template.New("database").Funcs(funcMap).Parse(database))
var nats = flag.String("nats", "nats", "path for nats folder")

func main() {
	databases := []string{"alert", "alert_option", "check", "client", "command", "group", "server", "upload", "user"}

	for _, db := range databases {
		file, err := openFile(path.Join(*nats, db+"s.go"))
		if err != nil {
			panic(err.Error())
		}

		data := struct {
			Name string
		}{db}

		err = databaseTemplate.Execute(file, data)
		file.Close()
		if err != nil {
			panic(err.Error())
		}
	}
}

func openFile(path string) (*os.File, error) {
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return nil, err
		}
		return file, err
	} else {
		err = os.Remove(path)
		if err != nil {
			return nil, err
		}

		var file, err = os.Create(path)
		if err != nil {
			return nil, err
		}
		return file, err
	}
}

var database = `package nats

import (
	"time"

	"github.com/keiwi/utils/models"
	"github.com/nats-io/go-nats"
	"gopkg.in/mgo.v2/bson"
)

func Delete{{CamelCase .Name}}(state *nats.Conn, data []byte) error {
	return state.Publish("{{.Name}}s.delete.send", data)
}

func Find{{CamelCase .Name}}(state *nats.Conn, data []byte) ([]models.{{CamelCase .Name}}, error) {
	// msg, err := state.Request("{{.Name}}s.retrieve.find", data, time.Duration(viper.GetInt("nats.delay"))*time.Second)
	msg, err := state.Request("{{.Name}}s.retrieve.find", data, time.Duration(10)*time.Second)
	if err != nil {
		return nil, err
	}

	var {{LowerCamelCase .Name}}s []models.{{CamelCase .Name}}
	err = bson.UnmarshalJSON(msg.Data, &{{LowerCamelCase .Name}}s)
	if err != nil {
		return nil, err
	}
	return {{LowerCamelCase .Name}}s, nil
}

func Has{{CamelCase .Name}}(state *nats.Conn, data []byte) (bool, error) {
	// msg, err := state.Request("{{.Name}}s.retrieve.has", data, time.Duration(viper.GetInt("nats.delay"))*time.Second)
	msg, err := state.Request("{{.Name}}s.retrieve.has", data, time.Duration(10)*time.Second)
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

func Create{{CamelCase .Name}}(state *nats.Conn, data []byte) error {
	return state.Publish("{{.Name}}s.create.send", data)
}

func Update{{CamelCase .Name}}(state *nats.Conn, data []byte) error {
	return state.Publish("{{.Name}}s.update.send", data)
}
`
