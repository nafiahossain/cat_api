package main

import (
	_ "cat_api/routers"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/core/config"
	"log"
)

func main() {
	// Load configuration
	conf, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Retrieve the API key from the config
	apiKey, err := conf.String("cat_api_key")
	if err != nil {
		log.Fatalf("Failed to get API key: %v", err)
	}
	
	if apiKey == "" {
		log.Fatalf("API key is missing in the configuration")
	}

	// Set Beego's run mode
	web.BConfig.RunMode = "prod" // or "dev", depending on your needs

	// Start the Beego server
	web.Run()
}
