package interfaces

import "leanmeal/api/models"

type ICustomerService interface {
	Add(name string, age int) bool
	Get() map[string]models.Player
	Remove(name string)
}
