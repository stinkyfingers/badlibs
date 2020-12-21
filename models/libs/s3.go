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

const bucket = "badlibs"
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
	j, err := json.Marshal(l)
	if err != nil {
		return err
	}
	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(l.ID),
		Body:   bytes.NewReader(j),
	})
	return err
}

func (l *Lib) Delete() error {
	sess, err := Session(region)
	if err != nil {
		return err
	}
	_, err = s3.New(sess).DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(l.ID),
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
		Key:    aws.String(l.ID),
	})
	if err != nil {
		return err
	}
	if resp == nil {
		return ErrNilObject
	}
	err = json.NewDecoder(resp.Body).Decode(l)
	return err
}

func (l *Lib) Find() ([]Lib, error) {
	sess, err := Session(region)
	if err != nil {
		return nil, err
	}
	resp, err := s3.New(sess).ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, ErrNilObject
	}
	// TODO pagination

	var libs []Lib
	for _, obj := range resp.Contents {
		resp, err := s3.New(sess).GetObject(&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(*obj.Key),
		})
		if err != nil {
			return nil, err
		}
		var l Lib
		err = json.NewDecoder(resp.Body).Decode(&l)
		if err != nil {
			return nil, err
		}
		libs = append(libs, l)
	}
	return libs, err
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
