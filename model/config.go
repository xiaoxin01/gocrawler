package model

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	// Webs webs to crawl
	Webs []Web
	// Db db connection info
	Db DbConnection
	// AlertKey alert key to send alert, http://sc.ftqq.com/
	AlertKey string
	// AlertType alert type
	AlertType string
)

// DbConnection connection info
type DbConnection struct {
	Connection string
	Database   string
	Collection string
}

// InitConfig initial config
func InitConfig(configPath string) {
	viper.AddConfigPath(configPath)
	viper.SetConfigName("webs")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	webs := viper.GetStringMap("webs")
	for path := range webs {
		var web Web
		if err := viper.UnmarshalKey("webs."+path, &web); err != nil {
			panic(err)
		}
		Webs = append(Webs, web)
	}

	err = viper.UnmarshalKey("db", &Db)
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("read db config err"))
	}

	AlertKey = viper.GetString("alertKey")
	AlertType = viper.GetString("alertType")
}
