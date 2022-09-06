package main

import (
	"os"

	"github.com/broem/live-music-archiver-api/server/api"
	_ "github.com/joho/godotenv/autoload"
	"gopkg.in/yaml.v2"
)

func main() {
	cfg := getConfig()
	api.NewApi(cfg)
}

func getConfig() *api.Config {
	f, err := os.Open("config.yml")
	if err != nil {
		return nil
	}
	defer f.Close()

	var cfg api.Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil
	}
	
	return &cfg
}
