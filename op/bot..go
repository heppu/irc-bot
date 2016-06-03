package ops

import (
	"encoding/binary"
	"fmt"
	"log"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/heppu/jun/client"
	"github.com/sorcix/irc"
)

const (
	BUCKET_PREFIX    = "ops_"
	COMMAND_1_PREFIX = "!op "
	COMMAND_2_PREFIX = ".op "
)

type Bot struct {
	irc *client.Client
	db  *bolt.DB
}

type UrlInfo struct {
	OpCount int
}

func NewBot(ircClient *client.Client, db *bolt.DB) *Bot {
	bot := &Bot{
		irc: ircClient,
		db:  db,
	}
	bot.addCallbacks()
	return bot
}

func (b *Bot) addCallbacks() {
	b.irc.AddCallback("PRIVMSG", func(message *irc.Message) {
		// Handle each meaasage on own goroutine
		go b.handlePrivMessage(message)
	})
}

func (b *Bot) handlePrivMessage(message *irc.Message) {
	// Check that this message was ment to our bot
	if !strings.HasPrefix(message.Trailing, COMMAND_1_PREFIX) &&
		!strings.HasPrefix(message.Trailing, COMMAND_2_PREFIX) {
		return
	}

	channel := message.Params[0]
	if channel == b.irc.Nickname {
		return
	}

	arr := strings.Split(message.Trailing, " ")
	filter(&arr, message.Name)

	if len(arr) < 2 {
		return
	}
	modes := "+"
	nicks := ""
	for _, nick := range arr[1:] {
		modes += "o"
		nicks += nick + " "
	}

	msg := fmt.Sprintf("/MODE %s %s", modes, nicks)
	b.irc.Raw(msg)
}

func filter(xs *[]string, me string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *xs {
		if !found[x] && x != me {
			found[x] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}
