package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/boltdb/bolt"
	"github.com/heppu/irc-bot/ops"
	"github.com/heppu/irc-bot/wanha"
	"github.com/heppu/jun/client"
)

const (
	DB_NAME    = "ircbot.db"
	IRC_SERVER = "irc.ca.ircnet.net:6667"
	BOT_NAME   = "bot-asd-123"
)

var channels = []string{"#gobot"}

func main() {
	db, err := bolt.Open(DB_NAME, 0600, nil)
	if err != nil {
		log.Fatalln(err)
	}

	delay := 500 * time.Microsecond
	c := client.New(
		IRC_SERVER,
		BOT_NAME,
		channels,
		nil,
		&delay,
	)

	// Create wanha bots
	wanha.NewBot(c, db)
	ops.NewBot(c, db)

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
