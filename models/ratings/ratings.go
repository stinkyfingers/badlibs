package ratings

import (
	"github.com/stinkyfingers/badlibs/helpers/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Rating struct {
	ID    bson.ObjectId `bson:"_id" json:"id"`
	Value string        `bson:"value" json:"value"`
}

func (r *Rating) FindMatch() ([]Rating, error) {
	var err error
	var rs []Rating
	session, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return rs, err
	}
	defer session.Close()
	c := session.DB(database.MongoDatabase()).C("ratings")

	querymap := make(map[string]interface{})
	if r.ID != "" {
		querymap["_id"] = r.ID
	}
	if r.Value != "" {
		querymap["value"] = r.Value
	}
	err = c.Find(querymap).All(&rs)
	return rs, err
}
