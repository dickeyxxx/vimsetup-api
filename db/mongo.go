package db

import (
	"labix.org/v2/mgo"
	"os"
)

func MongoSession() *mgo.Session {
	uri := os.Getenv("MONGOHQ_URL")
	if uri == "" {
		uri = "mongodb://localhost/vimsetup"
	}
	mongo, err := mgo.Dial(uri)
	if err != nil {
		panic(err)
	}

	return mongo
}
