package partsofspeech

import (
	"github.com/stinkyfingers/badlibs/helpers/database"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type PartOfSpeech struct {
	ID    bson.ObjectId `bson:"_id" json:"id"`
	Value string        `bson:"value" json:"value"`
	Code  string        `bson:"code" json:"code"`
}

func (p *PartOfSpeech) FindMatch() ([]PartOfSpeech, error) {
	var err error
	var ps []PartOfSpeech
	session, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return ps, err
	}
	defer session.Close()
	c := session.DB(database.MongoDatabase()).C("partsofspeech")

	querymap := make(map[string]interface{})
	if p.ID != "" {
		querymap["_id"] = p.ID
	}
	if p.Value != "" {
		querymap["value"] = p.Value
	}
	if p.Code != "" {
		querymap["code"] = p.Code
	}
	err = c.Find(querymap).All(&ps)
	return ps, err
}

func (p *PartOfSpeech) Create() error {
	session, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return err
	}
	defer session.Close()
	p.ID = bson.NewObjectId()
	return session.DB(database.MongoDatabase()).C("partsofspeech").Insert(p)
}
