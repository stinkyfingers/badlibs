package database

import (
	"flag"
	"fmt"
	"gopkg.in/mgo.v2"
	"os"
	"time"
)

var (
	EmptyDb = flag.String("clean", "", "bind empty database with structure defined")
)

func ConnectionString() string {

	if addr := os.Getenv("CLEARDB_DATABASE_URL"); addr != "" {
		return fmt.Sprint(addr)
	}

	if EmptyDb != nil && *EmptyDb != "" {
		return "root:@tcp(127.0.0.1:3306)/Kindling_Empty?parseTime=true&loc=America%2FChicago"
	}
	return "root:@tcp(127.0.0.1:3306)/Kindling?parseTime=true&loc=America%2FChicago"
}

func MongoConnectionString() *mgo.DialInfo {

	var (
		MongoDBHosts = os.Getenv("OPENSHIFT_MONGODB_DB_HOST")
		MongoPort    = os.Getenv("OPENSHIFT_MONGODB_DB_PORT")
		// AuthDatabase    = os.Getenv("MONGO_DB")
		AuthUserName    = os.Getenv("OPENSHIFT_MONGODB_DB_USERNAME")
		AuthPassword    = os.Getenv("OPENSHIFT_MONGODB_DB_PASSWORD")
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
