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
	COMMAND_1_PREFIX = "!ops"
	COMMAND_2_PREFIX = ".ops"
	INC_OPS          = "+o"
	DEC_OPS          = "-o"
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
	b.irc.AddCallback("MODE", func(message *irc.Message) {
		// Handle each meaasage on own goroutine
		go b.handleModeMessage(message)
	})

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
	removeDuplicates(&arr)

	if len(arr) == 1 {
		b.handleOpsPrint(channel, message.Name)
		return
	}

	for _, nick := range arr[1:] {
		if len(nick) > 0 {
			b.handleOpsPrint(channel, nick)
		}
	}
}

func (b *Bot) handleModeMessage(message *irc.Message) {
	if len(message.Params) < 3 {
		return
	}

	// Get channel name
	channel := message.Params[0]

	// Parse mode
	mode := message.Params[1]
	var f func(int64) int64

	// We count only ops
	if mode == INC_OPS {
		f = inc
	} else if mode == DEC_OPS {
		f = dec
	} else {
		return
	}

	nick := message.Params[2]
	b.handleOps(channel, nick, f)
}

func (b *Bot) handleOps(channel, nick string, f func(int64) int64) {
	log.Println("Handle nick:", nick)

	var i int64
	bucketName := BUCKET_PREFIX + channel
	bs := make([]byte, 4)

	err := b.db.Update(func(tx *bolt.Tx) (err error) {
		bucket := tx.Bucket([]byte(bucketName))

		// If bucket does not exist create one
		if bucket == nil {
			log.Println("No bucket found, creating:", bucketName)
			if bucket, err = tx.CreateBucket([]byte(bucketName)); err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}

			// Set ops count for nick
			i = f(i)
			binary.PutVarint(bs, i)
			if err = bucket.Put([]byte(nick), bs); err != nil {
				return
			}
			b.irc.Privmsg(channel, fmt.Sprintf("\x02%s ops count: %d", nick, i))
			return
		}

		v := bucket.Get([]byte(nick))
		if v != nil {
			i, _ = binary.Varint(v)
		}
		i = f(i)
		binary.PutVarint(bs, i)

		if err = bucket.Put([]byte(nick), bs); err != nil {
			return
		}

		b.irc.Privmsg(channel, fmt.Sprintf("\x02%s ops count: %d", nick, i))
		return
	})

	if err != nil {
		log.Println(err)
		return
	}
}

func (b *Bot) handleOpsPrint(channel, nick string) {
	var i int64
	bucketName := BUCKET_PREFIX + channel

	b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		log.Printf("'%s'  '%s'\n", bucketName, nick)

		if bucket == nil {
			b.irc.Privmsg(channel, "\x02No ops to show")
			return nil
		}

		v := bucket.Get([]byte(nick))
		if v == nil {
			return nil
		}
		i, _ = binary.Varint(v)
		b.irc.Privmsg(channel, fmt.Sprintf("%s ops count: %d", nick, i))
		return nil
	})
}

func dec(i int64) int64 {
	return i - 1
}

func inc(i int64) int64 {
	return i + 1
}

func removeDuplicates(xs *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *xs {
		if !found[x] {
			found[x] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}
