package implementations

import (
	"log"
	"os"
)

type Initialization struct {
	Storage Storage
}

func (i *Initialization) Initialized() bool {
	sqlBytes, err := os.ReadFile("db.sql")
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
		return false
	}

	content := string(sqlBytes)

	return content != "1"
}

func (i *Initialization) Database() (bool, error) {
	sqlBytes, err := os.ReadFile("db.sql")
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
		return false, err
	}

	sqlContent := string(sqlBytes)
	i.Storage.Open()
	defer i.Storage.Close()
	result := i.Storage.Exec(sqlContent, []interface{}{})

	file, err := os.Create("done")
	if err != nil {
		return false, err
	}

	file.Write([]byte{1})
	file.Close()

	return result, nil
}

func (i *Initialization) Seed() (bool, error) {
	sqlBytes, err := os.ReadFile("init.sql")
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
		return false, err
	}

	sqlContent := string(sqlBytes)
	i.Storage.Open()
	defer i.Storage.Close()
	result := i.Storage.Exec(sqlContent, []interface{}{})

	return result, nil
}
