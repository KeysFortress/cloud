package implementations

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	values map[string]interface{}
}

func (config *Configuration) Load() bool {
	config.values = make(map[string]interface{})

	file, err := os.Open("settings.json")
	if err != nil {
		panic("Configuration file doesn't exist.")
	}

	// Decode the JSON file into the map
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config.values); err != nil {
		fmt.Println("Error decoding JSON:", err)

	}

	defer file.Close()

	// Print the map
	fmt.Println(config.values)
	return true
}

func (config *Configuration) GetKey(name string) interface{} {
	result := config.values[name]
	return result
}
