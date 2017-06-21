package storage

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Bank struct {
	Name  string
	Token string
}

type Person struct {
	Email string
	Banks []Bank
}


type Storage interface {
	GetPerson(string) (*Person, error)
}

type FakeStorage struct {
	Tokens []string
}

func (f *FakeStorage) GetPerson(email string) (*Person, error) {
	p := &Person{
		Email: email,
		Banks: []Bank{Bank{"bank1", f.Tokens[0]}},
	}
	return p, nil
}

type MongoStorage struct {
	Session *mgo.Session
}

func (s *MongoStorage) GetPerson(email string) (*Person, error) {
	c := s.Session.DB("test").C("people")

	result := new(Person)
	err := c.Find(bson.M{"name": email}).One(result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

