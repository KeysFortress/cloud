package repositories

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"leanmeal/api/interfaces"
)

type AccessKeysRepository struct {
	ConnectionString string
	Storage          interfaces.Storage
}

func (accessKeys *AccessKeysRepository) Add(id *uuid.UUID, key *string) (uuid.UUID, error) {
	query := `
					INSERT INTO public.associated_account_access_keys(
					account_id, key, created_at)
					VALUES ($1, $2, $3)
					RETURNING id
			`

	var createdId uuid.UUID
	queryResult := accessKeys.Storage.Add(&query, &[]interface{}{&id, &key, time.Now().UTC()})

	err := queryResult.Scan(&createdId)
	if err != nil {
		fmt.Printf("Failed to add access key for account %v", id)
		fmt.Println(err)
		return uuid.UUID{}, err
	}

	return createdId, nil
}

func (accesKeys *AccessKeysRepository) GetAccountKeys(id uuid.UUID) []string {
	query := `
				SELECT access_key as "keyVal" FROM public.associated_account_access_keys
				WHERE account_id = $1
			`

	rows := accesKeys.Storage.Where(query, []interface{}{id})

	var keys []string
	for rows.Next() {
		var keyVal string
		err := rows.Scan(&keyVal)
		if err != nil {
			fmt.Println(err)
		}

		keys = append(keys, keyVal)
	}

	return keys
}
