package libs

import (
	"encoding/json"
	"os"
	"time"

	libs "github.com/stinkyfingers/badlibs/models"

	"github.com/google/uuid"
)

type FileStorage struct {
	file string
}

type DBMap map[string]libs.Lib

func NewFileStorage(file string) (*FileStorage, error) {
	return &FileStorage{
		file: file,
	}, nil
}

func (s *FileStorage) Create(l *libs.Lib) (*libs.Lib, error) {
	ti := time.Now()
	l.Created = &ti
	l.ID = uuid.New().String()
	return s.Update(l)
}

func (s *FileStorage) Update(l *libs.Lib) (*libs.Lib, error) {
	f, err := os.Open(s.file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var libs DBMap
	err = json.NewDecoder(f).Decode(&libs)
	if err != nil {
		return nil, err
	}
	libs[l.ID] = *l
	f.Close()
	f, err = os.Create(s.file)
	if err != nil {
		return nil, err
	}
	err = json.NewEncoder(f).Encode(libs)
	return l, err
}

func (s *FileStorage) Delete(id string) error {
	f, err := os.Open(s.file)
	if err != nil {
		return err
	}
	defer f.Close()
	var libs DBMap
	err = json.NewDecoder(f).Decode(&libs)
	if err != nil {
		return err
	}
	delete(libs, id)
	f.Close()
	f, err = os.Create(s.file)
	if err != nil {
		return err
	}
	return json.NewEncoder(f).Encode(libs)
}

func (s *FileStorage) Get(id string) (*libs.Lib, error) {
	f, err := os.Open(s.file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var libs DBMap
	err = json.NewDecoder(f).Decode(&libs)
	if err != nil {
		return nil, err
	}
	lib := libs[id]
	return &lib, nil
}

func (s *FileStorage) All(filter *libs.Lib) ([]libs.Lib, error) {
	f, err := os.Open(s.file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var dbMap DBMap
	err = json.NewDecoder(f).Decode(&dbMap)
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
