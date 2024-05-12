package repositories

import (
	"leanmeal/api/interfaces"
	"leanmeal/api/models"
)

type PasswordRepository struct {
	Storage interfaces.Storage
}

func (pr *PasswordRepository) All() ([]models.Password, error) {
	query := `
		SELECT id,name, website, LENGTH(content) as password_lenght, created_at, updated_at
		FROM public.account_passwords
		GROUP BY id
	`

	queryResult := pr.Storage.Where(query, []interface{}{})

	var passwords []models.Password
	for queryResult.Next() {
		var password models.Password
		err := queryResult.Scan(&password.Id, &password.Email, &password.Website, &password.Password,
			&password.CreatedAt, &password.UpdatedAt)

		if err != nil {
			return []models.Password{}, err
		}

		passwords = append(passwords, password)
	}

	return passwords, nil
}
