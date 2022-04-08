package libs

import (
	"time"
)

type Lib struct {
	ID      string     `json:"_id"`
	Text    string     `json:"text"`
	Title   string     `json:"title"`
	Rating  string     `json:"rating"`
	Created *time.Time `json:"created"`
	User    string     `json:"user"`
}

type LibStorer interface {
	Get(id string) (*Lib, error)
	All() ([]Lib, error)
	Delete(id string) error
	Update(l *Lib) (*Lib, error)
	Create(l *Lib) (*Lib, error)
}

