package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	Db map[string]string
	Imdb map[string]string
	Api map[string]int	
}

func GetConfig() Configuration {
	file, err := os.Open("config/config.json")
	if err != nil {
		fmt.Println("error:", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	return configuration
}
