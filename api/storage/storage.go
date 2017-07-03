package storage

import (
  "fmt"
  "log"

	"gopkg.in/mgo.v2"
)

type Account struct {
	Name  string
	Token string
}

type Person struct {
	Email string
	Accounts []Account
}


type Storage interface {
	GetPerson(string) (*Person, error)
  Exists(string) (bool, error)

  // Returns true if the person was created
  CreatePerson(string) (bool, error)
}

type FakeStorage struct {
	Tokens []string
}

func (f *FakeStorage) GetPerson(email string) (*Person, error) {
	p := &Person{
		Email: email,
		Accounts: []Account{Account{"bank1", f.Tokens[0]}},
	}
	return p, nil
}

func (f *FakeStorage) Exists(email string) (bool, error) {
  return true, nil
}

func (f *FakeStorage) CreatePerson(email string) (bool, error) {
  return false, nil
}

type MongoStorage struct {
	Session *mgo.Session
}

func (s *MongoStorage) GetPerson(email string) (*Person, error) {
	c := s.Session.DB("test").C("people")

  person := Person{Email: email}
	err := c.Find(&person).One(&person)

	if err != nil {
		return nil, err
	}

	return &person, nil
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


func (s *MongoStorage) CreatePerson(email string) (bool, error) {
  exists, err := s.Exists(email)

  if err != nil {
    return false, err
  }

  if exists {
    return false, nil
  }

  log.Printf("CreatePerson")

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
