package storage

import (
  "fmt"
  "log"

	"gopkg.in/mgo.v2"
)

type Account struct {
  Name  string `json:"name"`
  Token string `json:"token"`
}

type Person struct {
  Email string       `bson:"email"`
	Accounts []Account `bson:"accounts,omitempty"`
}


type Storage interface {
	Get(string) (*Person, error)
  Exists(string) (bool, error)

  // Returns true if created, false if it already existed
  Create(string) (bool, error)

  // Update person
  Update(*Person) error
}

type FakeStorage struct {
	Tokens []string
}

func (f *FakeStorage) Get(email string) (*Person, error) {
	p := &Person{
		Email: email,
		Accounts: []Account{Account{"bank1", f.Tokens[0]}},
	}
	return p, nil
}

func (f *FakeStorage) Exists(email string) (bool, error) {
  return true, nil
}

func (f *FakeStorage) Create(email string) (bool, error) {
  return false, nil
}

func (f *FakeStorage) Update(*Person) error {
  return nil
}

type MongoStorage struct {
	Session *mgo.Session
}

func (s *MongoStorage) Get(email string) (*Person, error) {
	c := s.Session.DB("test").C("people")

  person := Person{Email: email}
  result := new(Person)
	err := c.Find(person).One(result)

	if err != nil {
		return nil, err
	}

	return result, nil
}


func (s *MongoStorage) Exists(email string) (bool, error) {
	c := s.Session.DB("test").C("people")
  person := Person{Email: email}
	count, err := c.Find(&person).Count()

	if err != nil {
		return false, err
	}

  if count > 1 {
    return false, fmt.Errorf("multiple people found for %s", email)
  }

	return count == 1, nil
}

func (s *MongoStorage) Update(p *Person) error {
	c := s.Session.DB("test").C("people")
  selector := Person{Email: p.Email}
  return c.Update(&selector, p)
}

func (s *MongoStorage) Create(email string) (bool, error) {
  exists, err := s.Exists(email)

  if err != nil {
    return false, err
  }

  if exists {
    return false, nil
  }

  log.Printf("Create")

	c := s.Session.DB("test").C("people")
	p := &Person{
		Email: email,
	}
  err = c.Insert(p)

  if err != nil {
    return false, err
  }

	return true, nil
}
