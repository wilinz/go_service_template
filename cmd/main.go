package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"server_template/common"
	"server_template/db"
	"server_template/route"
	"server_template/service/proxy"
	"server_template/tools"
)

type ServiceConfig struct {
	MysqlConfig       db.MysqlConfig          `json:"mysql_config"`
	RedisConfig       db.RedisConfig          `json:"redis_config"`
	EmailConfig       tools.EmailConfig       `json:"email_config"`
	CustomProxyConfig proxy.CustomProxyConfig `json:"custom_proxy_config"`
}

func main() {
	log.SetFlags(log.Llongfile | log.LstdFlags)

	// Define command line flags
	configFile := flag.String("c", "", "Path to the service configuration file")
	genTemplate := flag.Bool("g", false, "Generate template service configuration file")

	// Parse command line flags
	flag.Parse()

	if *genTemplate {
		// Generate template configuration file
		config := ServiceConfig{
			MysqlConfig: db.MysqlConfig{
				Host:     "localhost",
				Port:     3306,
				Username: "root",
				Password: "",
				Database: "test",
			},
			RedisConfig: db.RedisConfig{
				Host: "localhost",
				Port: 6379,
			},
			EmailConfig: tools.EmailConfig{
				SMTPHost:     "smtp.example.com",
				SMTPPort:     587,
				SMTPUsername: "your_email@example.com",
				SMTPPassword: "your_email_password",
				FromAddress:  "your_email@example.com",
			},
			CustomProxyConfig: proxy.CustomProxyConfig{
				ProxyServer: "localhost",
				Key:         "",
			},
		}

		// Marshal the configuration to JSON
		configJSON, err := json.MarshalIndent(config, "", "    ")
		if err != nil {
			log.Fatal(err)
		}

		// Write the configuration to a file
		err = ioutil.WriteFile("service_config_temp.json", configJSON, 0644)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Template configuration file generated at 'service_config_temp.json'")
		return
	}

	// If no configuration file specified, try to read default configuration file in current directory
	if *configFile == "" {
		if _, err := os.Stat("service_config.json"); err == nil {
			*configFile = "service_config.json"
		} else if _, err := os.Stat("service_config.json5"); err == nil {
			*configFile = "service_config.json5"
		} else {
			log.Fatal("No configuration file specified. Use the '-c' flag to specify the path to the configuration file or place a default configuration file 'service_config.json' or 'service_config.json5' in the current directory.")
		}
	}

	// Read the configuration file
	configJSON, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the configuration JSON into a ServiceConfig struct
	var config ServiceConfig
	err = json.Unmarshal(configJSON, &config)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the service using the configuration
	db.InitMySql(config.MysqlConfig)
	db.InitRedis(config.RedisConfig)
	tools.InitEmail(config.EmailConfig)
	proxy.InitCustomProxy(config.CustomProxyConfig)
	common.InitSessions()
	route.Run()
}
