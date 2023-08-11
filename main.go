package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/titanous/json5"
	"io/ioutil"
	"log"
	"os"
	"server_template/common"
	"server_template/db"
	"server_template/route"
	"server_template/service/account"
	"server_template/service/proxy"
	"server_template/tools"
)

type ServiceConfig struct {
	AppName           string                  `json:"app_name"`
	Host              string                  `json:"host"`
	Port              int                     `json:"port"`
	BaseUrl           string                  `json:"base_url"`
	UploadDir         string                  `json:"upload_dir"`
	MysqlConfig       db.MysqlConfig          `json:"mysql_config"`
	RedisConfig       db.RedisConfig          `json:"redis_config"`
	EmailConfig       tools.EmailConfig       `json:"email_config"`
	CustomProxyConfig proxy.CustomProxyConfig `json:"custom_proxy_config"`
}

// @title server_template API
// @version 1.0
// @description server_template API
// @termsOfService https://github.com/wilinz/server_template

// @contact.name API Support
// @contact.url https://github.com/wilinz/server_template
// @contact.email weilizan71@outlook.com

// @in header

// 下面注释按照项目实际的地址和路径进行设置
// @host 127.0.0.1:10010
// @BasePath /
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
			AppName: "GoodApp",
			Host:    "0.0.0.0",
			Port:    10010,
			BaseUrl: "",
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
				Timeout:     300,
			},
		}

		// Marshal the configuration to JSON
		configJSON, err := json.MarshalIndent(config, "", "    ")
		if err != nil {
			log.Fatal(err)
		}

		// Write the configuration to a file
		err = ioutil.WriteFile("service_config_temp.json5", configJSON, 0644)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Template configuration file generated at 'service_config_temp.json5'")
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
	err = json5.Unmarshal(configJSON, &config)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the service using the configuration
	common.InitUploadSavePath(config.UploadDir)
	db.InitMySql(config.MysqlConfig)
	db.InitRedis(config.RedisConfig)
	account.InitAppName(config.AppName)
	tools.InitEmail(config.EmailConfig)
	common.InitSessions()
	proxy.InitCustomProxy(config.CustomProxyConfig)
	route.Run(config.Host, config.Port)
}
