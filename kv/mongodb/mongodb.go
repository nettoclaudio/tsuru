package mongodb

import (
	"fmt"

	"github.com/tsuru/tsuru/app/bind"
	"github.com/tsuru/tsuru/db"
	"gopkg.in/mgo.v2/bson"
)

type mongoDBKeyValueStorager struct {
	appName string
}

type envdb struct {
	Env map[string]bind.EnvVar `bson:"env"`
}

func (es *mongoDBKeyValueStorager) Get(envs ...string) (map[string]bind.EnvVar, error) {
	var envv envdb

	conn, err := db.Conn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	err = conn.Apps().Find(bson.M{"name": es.appName}).Select(bson.M{"env": true, "_id": false}).One(&envv)

	envsDB := envv.Env

	if err != nil {
		return nil, err
	}

	if len(envs) > 0 {

		result := make(map[string]bind.EnvVar, 0)

		for _, env := range envs {
			fmt.Printf("%#v\n", env)

			if v, ok := envsDB[env]; ok {
				fmt.Printf("%#v\n", v)
				result[env] = v

				fmt.Printf("%#v\n", result)
			}
		}

		return result, nil
	}

	return envsDB, nil
}

func (es *mongoDBKeyValueStorager) Set(envs []bind.EnvVar) error {
	return nil
}

func (es *mongoDBKeyValueStorager) Unset(envs ...string) error {
	return nil
}

func NewMongoDBKeyValueStorager(appName string) *mongoDBKeyValueStorager {
	return &mongoDBKeyValueStorager{appName: appName}
}
