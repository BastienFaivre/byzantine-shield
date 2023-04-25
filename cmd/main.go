package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BastienFaivre/byzantine-shield/internal/communication"
	"github.com/BastienFaivre/byzantine-shield/internal/config"
)

var (
	configPath string
	addr       string
	port       int
)

func init() {
	flag.StringVar(&configPath, "config", "", "Path to the configuration file")
	flag.StringVar(&addr, "addr", "127.0.0.1", "Listen address")
	flag.IntVar(&port, "port", 8080, "Listen port")
}

func main() {
	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("BYZANTINE_SHIELD_CONFIG")
		if configPath == "" {
			log.Fatal("Please specify a configuration file via --config or the BYZANTINE_SHIELD_CONFIG environment variable")
		}
	}

	config, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	proxy := communication.NewProxy(*config)
	http.HandleFunc("/", proxy.HandleRequest)

	log.Printf("Listening on %s:%d", addr, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), nil))
}
