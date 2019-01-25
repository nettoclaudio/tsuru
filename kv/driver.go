package kv

import (
	"github.com/tsuru/tsuru/app/bind"
	"github.com/tsuru/tsuru/kv/mongodb"
)

type KeyValueStorager interface {
	Get(...string) (map[string]bind.EnvVar, error)
	Set([]bind.EnvVar) error
	Unset(...string) error
}

func GetDefaultKeyValueStorager(appName string) KeyValueStorager {
	return mongodb.NewMongoDBKeyValueStorager(appName)
}
