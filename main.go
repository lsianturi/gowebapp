package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/lsianturi/gowebapp/daemon"
)

var assetPath string

func processFlag() *daemon.Config {
	cfg := &daemon.Config{}
	flag.StringVar(&cfg.ListenSpec, "listen", "localhost:3000", "HTTP listen spec")
	flag.StringVar(&cfg.Db.ConnectString, "db-connect", "gowebapp:sipuserix@tcp(10.15.2.115:3306)/demo", "DB Connect String")
	flag.StringVar(&assetPath, "assets-path", "assets", "Path to assets dir")

	flag.Parse()
	return cfg
}

func setupHttpAssets(cfg *daemon.Config) {
	log.Printf("Assets served from %q.", assetPath)
	cfg.UI.Assets = http.Dir(assetPath)
}

func main() {
	cfg := processFlag()

	setupHttpAssets(cfg)

	if err := daemon.Run(cfg); err != nil {
		log.Printf("Error in main(): %v", err)
	}
}
