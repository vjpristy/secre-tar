package main

import (
	"log"
	"net/url"

	"github.com/vjpristy/secre-tar/internal/config"
	"github.com/vjpristy/secre-tar/internal/gui"
	"github.com/vjpristy/secre-tar/internal/network"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	u := url.URL{Scheme: "ws", Host: cfg.ServerAddress, Path: "/ws"}
	conn, _, err := network.Dial(u.String())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	secretaryGUI := gui.NewSecretaryGUI(conn)
	secretaryGUI.Run()
}
