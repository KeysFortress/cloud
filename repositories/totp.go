package repositories

import "leanmeal/api/interfaces"

type TotpRepository struct {
	Storage interfaces.Storage
}
