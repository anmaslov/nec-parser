package store

import (
	"fmt"
	"github.com/anmaslov/nec-parser/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type MongoStore struct {
	Session *mgo.Session
	Db      string
}

// NewMongo новое подключение к монго
func NewMongo(server *config.Db) (*MongoStore, error) {

	conn := &mgo.DialInfo{
		Addrs:    []string{server.Host},
		Timeout:  5 * time.Second,
		Database: server.Dbname,
		Username: server.Username,
		Password: server.Password,
	}

	session, err := mgo.DialWithInfo(conn)
	if err != nil {
		return nil, fmt.Errorf("mongodb error: %s", err)
	}

	return &MongoStore{
		Session: session,
		Db:      server.Dbname,
	}, nil
}

// InsertCall Вставить в бд звонок
func (ms *MongoStore) InsertCall(call *CallInfo) error {
	c := ms.Session.DB(ms.Db).C("calls")

	err := c.Insert(call)
	if err != nil {
		return fmt.Errorf("error when trying write to mongoDB: %s", err)
	}

	return nil
}

// GetPhones Получить список телефонных станций
func (ms *MongoStore) GetPhones() ([]Phones, error) {
	c := ms.Session.DB(ms.Db).C("phones")

	query := bson.M{}
	query["enabled"] = bson.M{"$eq": true}

	phones := make([]Phones, 0)
	err := c.Find(query).Sort("-_id").All(&phones)
	if err != nil {
		return nil, fmt.Errorf("mongodb error: %s", err)
	}

	return phones, nil
}
