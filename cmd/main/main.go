package main

import (
	"flag"
	"log"

	"github.com/furrygem/nocut-api/internal/apiserver"
	"github.com/furrygem/nocut-api/pkg/logging"
)

var configFilePathYAML string

func main() {
	l := logging.GetLogger()
	l.Infoln("hello")
	flag.Parse()
	config := apiserver.NewConfig()
	if configFilePathYAML == "" {
		log.Println("Configuration file not specified, dropping to defaults.")
	} else {
		log.Println("Configuring via YAML file")
		err := config.FromYAML(configFilePathYAML)
		if err != nil {
			log.Fatal(err)
		}
	}

	server := apiserver.New(config)
	log.Printf("Listening on %s:%d", config.BindAddr, config.BindPort)
	if err := server.Start(config); err != nil {
		log.Fatal(err)
	}

}

func init() {
	flag.StringVar(&configFilePathYAML, "config-yaml", "", "Path to YAML config")
}
