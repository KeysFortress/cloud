package repositories

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"leanmeal/api/dtos"
	"leanmeal/api/interfaces"
	"leanmeal/api/models"
)

type Accounts struct {
	ConnectionString string
	Storage          interfaces.Storage
}

func (accountService *Accounts) UserExists(email string) (models.Account, error) {
	var account models.Account
	data := accountService.Storage.Single("select id, email from public.accounts where email = $1", []interface{}{email})

	err := data.Scan(&account.Id, &account.Email)
	if err != nil {
		fmt.Printf("Failed to fetch account with email %v", email)
		fmt.Println(err)
		return account, err
	}

	fmt.Println(&account)
	return account, nil
}

func (accountService *Accounts) GetById(id uuid.UUID) (models.Account, error) {
	var account models.Account
	data := accountService.Storage.Single("select id, email from public.accounts where id = $1", []interface{}{&id})

	err := data.Scan(&account.Id, &account.Email)
	if err != nil {
		fmt.Printf("Failed to fetch account with id %v", id)
		fmt.Println(err)
		return account, err
	}

	fmt.Println(&account)
	return account, nil
}

func (accountService *Accounts) CreateAccount(newAccount *dtos.CreateAccountRequest) (uuid.UUID, error) {
	_, err := accountService.UserExists(newAccount.Email)

	if err != nil {
		return uuid.UUID{}, nil
	}

	sql := `
	INSERT INTO public.accounts(
		 email, name, created_at, enabled)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	queryResult := accountService.Storage.Add(&sql, &[]interface{}{
		newAccount.Email,
		newAccount.Name,
		time.Now().UTC(),
		true,
	})
	var createdId uuid.UUID

	err = queryResult.Scan(&createdId)

	if err != nil {
		fmt.Printf("Failed to add access key for account %v", newAccount.Email)
		fmt.Println(err)
	}

	return createdId, nil
}

func (accountService *Accounts) Get() []models.Account {

	rows := accountService.Storage.Where("SELECT * from public.accounts", []interface{}{})

	var accounts []models.Account

	for rows.Next() {
		var account models.Account
		rows.Scan(&account.Id, &account.Email, &account.Name, &account.CreatedAt, &account.Enabled)
		accounts = append(accounts, account)
		fmt.Println(account)
	}

	return accounts
}
