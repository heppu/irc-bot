package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/boltdb/bolt"
	"github.com/heppu/irc-bot/ops"
	"github.com/heppu/irc-bot/wanha"
	"github.com/heppu/jun/client"
	"github.com/naoina/toml"
)

type Config struct {
	Name     string
	Channels []string
	Server   string
	Database string
	Delay    uint
}

const CONF_FILE = "conf.toml"

var conf Config

// Init configurations
func init() {
	file, err := os.Open(CONF_FILE)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	if err := toml.Unmarshal(buf, &conf); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Open boltdb
	db, err := bolt.Open(conf.Database, 0600, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Initialize new irc client
	delay := time.Duration(conf.Delay) * time.Millisecond
	c := client.New(
		conf.Server,
		conf.Name,
		conf.Channels,
		nil,
		&delay,
	)

	// Create bots
	wanha.NewBot(c, db)
	ops.NewBot(c, db)

	// Connect to server
	if err = c.Connect(); err != nil {
		log.Fatal(err)
	}

	// Shutdown if irc client fails
	go func(c *client.Client, db *bolt.DB) {
		err := <-c.Error

		c.Disconnect()
		db.Close()

		log.Fatalln(err)
	}(c, db)

	// Graceful shutdown for Ctrl+C
	go func(c *client.Client, db *bolt.DB) {
		kill := make(chan os.Signal, 1)
		signal.Notify(kill, os.Interrupt)
		<-kill

		c.Disconnect()
		db.Close()

		os.Exit(0)
	}(c, db)

	<-c.Quit
}
