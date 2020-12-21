package database

import (
	"os"
	"time"

	"gopkg.in/mgo.v2"
)

func MongoConnectionString() *mgo.DialInfo {

	var (
		MongoDBHosts = os.Getenv("MONGODB_DB_HOST")
		MongoPort    = os.Getenv("MONGODB_DB_PORT")
		// AuthDatabase    = os.Getenv("MONGO_DB")
		AuthUserName    = os.Getenv("MONGODB_DB_USERNAME")
		AuthPassword    = os.Getenv("MONGODB_DB_PASSWORD")
		mongoDBDialInfo mgo.DialInfo
	)

	if MongoDBHosts == "" {
		mongoDBDialInfo = mgo.DialInfo{
			Addrs:    []string{"127.0.0.1"},
			Timeout:  60 * time.Second,
			Database: "",
			Username: "",
			Password: "",
		}
	} else {
		addr := MongoDBHosts + ":" + MongoPort
		mongoDBDialInfo = mgo.DialInfo{
			Addrs:    []string{addr},
			Timeout:  60 * time.Second,
			Database: "badlibs",
			Username: AuthUserName,
			Password: AuthPassword,
		}
	}

	return &mongoDBDialInfo
}

func MongoDatabase() string {
	return "badlibs"
}
