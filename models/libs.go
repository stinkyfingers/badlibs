package libs

import (
	"time"
)

type Lib struct {
	ID      string     `json:"_id"`
	Text    string     `json:"text"`
	Title   string     `json:"title"`
	Rating  string     `json:"rating"`
	Rank    int        `json:"rank"`  // average rank: 1-10
	Ranks   int        `json:"ranks"` // number of times ranked
	Created *time.Time `json:"created"`
	User    User     `json:"user"`
}

type User struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

type LibStorer interface {
	Get(id string) (*Lib, error)
	All(filter *Lib) ([]Lib, error)
	Delete(id string) error
	Update(l *Lib) (*Lib, error)
	Create(l *Lib) (*Lib, error)
}
