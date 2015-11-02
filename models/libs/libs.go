package libs

import (
	"github.com/stinkyfingers/badlibs/helpers/database"
	// "github.com/stinkyfingers/badlibs/models/partsofspeech"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// "regexp"
	// "strings"
	"time"
)

type Lib struct {
	ID      bson.ObjectId `json:"_id" bson:"_id"`
	Text    string        `json:"text" bson:"text"`
	Title   string        `json:"title" bson:"title"`
	Rating  string        `json:"rating" bson:"rating"`
	Created *time.Time    `json:"created" bson:"created"`
}

func (l *Lib) Create() error {
	var err error
	session, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return err
	}
	defer session.Close()
	l.ID = bson.NewObjectId()

	c := session.DB(database.MongoDatabase()).C("libs")
	return c.Insert(l)
}

func (l *Lib) Update() error {
	var err error
	session, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return err
	}
	defer session.Close()

	c := session.DB(database.MongoDatabase()).C("libs")
	return c.UpdateId(l.ID, l)
}

func (l *Lib) Delete() error {
	var err error
	session, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return err
	}
	defer session.Close()

	c := session.DB(database.MongoDatabase()).C("libs")
	return c.Remove(l)

}

func (l *Lib) Get() error {
	var err error
	session, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return err
	}
	defer session.Close()
	c := session.DB(database.MongoDatabase()).C("libs")
	return c.FindId(l.ID).One(&l)
}

func (l *Lib) Find() ([]Lib, error) {
	var err error
	var ls []Lib
	session, err := mgo.DialWithInfo(database.MongoConnectionString())
	if err != nil {
		return ls, err
	}
	defer session.Close()
	c := session.DB(database.MongoDatabase()).C("libs")

	querymap := make(map[string]interface{})
	if l.ID != "" {
		querymap["_id"] = l.ID
	}
	if l.Text != "" {
		querymap["text"] = l.Text
	}
	if l.Title != "" {
		querymap["title"] = l.Title
	}
	if l.Rating != "" {
		querymap["rating"] = l.Rating
	}
	err = c.Find(querymap).All(&ls)
	if err != nil {
		return ls, err
	}
	return ls, err
}
