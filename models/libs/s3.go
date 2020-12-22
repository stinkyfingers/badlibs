package libs

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type Lib struct {
	ID      string     `json:"_id" bson:"_id"`
	Text    string     `json:"text" bson:"text"`
	Title   string     `json:"title" bson:"title"`
	Rating  string     `json:"rating" bson:"rating"`
	Created *time.Time `json:"created" bson:"created"`
}

type DBMap map[string]Lib

const bucket = "badlibs"
const key = "db.json"
const region = "us-west-1"

var (
	ErrNilObject = errors.New("nil object")
)

func (l *Lib) Create() error {
	l.ID = uuid.New().String()
	return l.Update()
}

func (l *Lib) Update() error {
	sess, err := Session(region)
	if err != nil {
		return err
	}

	resp, err := s3.New(sess).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	var libs DBMap
	err = json.NewDecoder(resp.Body).Decode(&libs)
	if err != nil {
		return err
	}
	libs[l.ID] = *l

	j, err := json.Marshal(libs)
	if err != nil {
		return err
	}
	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(j),
	})
	return err
}

func (l *Lib) Delete() error {
	sess, err := Session(region)
	if err != nil {
		return err
	}
	resp, err := s3.New(sess).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	if resp == nil {
		return ErrNilObject
	}
	var libs DBMap
	err = json.NewDecoder(resp.Body).Decode(&libs)
	if err != nil {
		return err
	}
	delete(libs, l.ID)
	j, err := json.Marshal(libs)
	if err != nil {
		return err
	}
	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(j),
	})
	return err
}

func (l *Lib) Get() error {
	sess, err := Session(region)
	if err != nil {
		return err
	}
	resp, err := s3.New(sess).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	if resp == nil {
		return ErrNilObject
	}
	var libs DBMap
	err = json.NewDecoder(resp.Body).Decode(&libs)
	if err != nil {
		return err
	}
	lib := libs[l.ID]
	l = &lib
	return nil
}

func (l *Lib) All() ([]Lib, error) {
	sess, err := Session(region)
	if err != nil {
		return nil, err
	}
	resp, err := s3.New(sess).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, ErrNilObject
	}
	var libs DBMap
	err = json.NewDecoder(resp.Body).Decode(&libs)
	if err != nil {
		return nil, err
	}
	var output []Lib
	for _, lib := range libs {
		output = append(output, lib)
	}
	return output, err
}

func Session(region string) (*session.Session, error) {
	options := session.Options{}
	if os.Getenv("PROFILE") != "" {
		options.Profile = os.Getenv("PROFILE")
	}

	sess, err := session.NewSessionWithOptions(options)
	if err != nil {
		return nil, err
	}

	sess.Config.WithRegion(region)

	return sess, nil
}

func AssureDBBucket() error {
	sess, err := Session(region)
	if err != nil {
		return err
	}
	_, err = s3.New(sess).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Body:   bytes.NewReader([]byte("{}")),
		})
	}
	return err
}
