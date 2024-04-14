package implementations

import (
	"leanmeal/api/models"
)

type CustomerService struct {
	People map[string]models.Player
}

func (cs *CustomerService) Add(name string, age int) bool {
	if cs.People == nil {
		cs.People = make(map[string]models.Player)
	}

	person := models.Player{
		Name: name,
		Age:  age,
	}

	cs.People[name] = person

	return true
}

func (cs *CustomerService) Get() map[string]models.Player {
	return cs.People
}

func (cs *CustomerService) Remove(name string) {
	delete(cs.People, name)
}
