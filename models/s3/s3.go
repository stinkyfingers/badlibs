package libs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	libs "github.com/stinkyfingers/badlibs/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/google/uuid"
)

type S3Storage struct {
	client s3iface.S3API
}

type DBMap map[string]libs.Lib

const bucket = "badlibs"
const key = "db.json"
const region = "us-west-1"

var (
	ErrNilObject = errors.New("nil object")
)

func NewS3Storage(profile string) (*S3Storage, error) {
	sess, err := Session(profile, region)
	if err != nil {
		return nil, err
	}
	client := s3.New(sess)
	err = AssureDBBucket(client)
	fmt.Println("BUCK", err, profile)
	if err != nil {
		return nil, err
	}
	return &S3Storage{
		client: client,
	}, nil
}

func (s *S3Storage) Create(l *libs.Lib) (*libs.Lib, error) {
	ti := time.Now()
	l.Created = &ti
	l.ID = uuid.New().String()
	return s.Update(l)
}

func (s *S3Storage) Update(l *libs.Lib) (*libs.Lib, error) {
	resp, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	var libs DBMap
	err = json.NewDecoder(resp.Body).Decode(&libs)
	if err != nil {
		return nil, err
	}
	libs[l.ID] = *l

	j, err := json.Marshal(libs)
	if err != nil {
		return nil, err
	}
	_, err = s.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(j),
	})
	return l, err
}

func (s *S3Storage) Delete(id string) error {
	resp, err := s.client.GetObject(&s3.GetObjectInput{
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
	delete(libs, id)
	j, err := json.Marshal(libs)
	if err != nil {
		return err
	}
	_, err = s.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(j),
	})
	return err
}

func (s *S3Storage) Get(id string) (*libs.Lib, error) {
	resp, err := s.client.GetObject(&s3.GetObjectInput{
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
	lib := libs[id]
	return &lib, nil
}

func (s *S3Storage) All(filter *libs.Lib) ([]libs.Lib, error) {
	resp, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, ErrNilObject
	}
	var dbMap DBMap
	err = json.NewDecoder(resp.Body).Decode(&dbMap)
	if err != nil {
		return nil, err
	}
	var output []libs.Lib
	for _, lib := range dbMap {
		if filter != nil {
			if filter.ID != "" && filter.ID != lib.ID ||
				filter.Title != "" && filter.Title != lib.Title ||
				filter.Rating != "" && filter.Rating != lib.Rating ||
				filter.User != "" && filter.User != lib.User ||
				filter.Created != nil && !filter.Created.IsZero() && filter.Created.After(*filter.Created) {
				continue
			}
		}
		output = append(output, lib)
	}
	return output, err
}

func Session(profile, region string) (*session.Session, error) {
	sess, err :=session.NewSession()
	if err != nil {
		return nil, err
	}

	sess.Config.WithRegion(region)
	return sess, nil
}

func AssureDBBucket(client *s3.S3) error {
	_, err := client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		_, err = client.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Body:   bytes.NewReader([]byte("{}")),
		})
	}
	return err
}