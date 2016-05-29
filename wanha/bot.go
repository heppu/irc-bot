package wanha

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/heppu/jun/client"
	"github.com/mvdan/xurls"
	"github.com/sorcix/irc"
)

const (
	BUCKET_PREFIX = "wanha_"
)

type Bot struct {
	irc *client.Client
	db  *bolt.DB
}

type UrlInfo struct {
	Message string
	Nick    string
	Time    int64
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
		go b.handleMessage(message)
	})
}

func (b *Bot) handleMessage(message *irc.Message) {
	// Parse urls from message and return if none were found
	urls := xurls.Relaxed.FindAllString(message.Trailing, -1)
	if len(urls) == 0 {
		return
	}

	// Get name for bucket which is channel or user's nick
	bucket := message.Params[0]
	if bucket == b.irc.Nickname {
		bucket = message.Name
	}

	// Handle each url
	for _, url := range urls {
		if strings.HasPrefix(url, "https://") {
			url = url[8:]
		} else if strings.HasPrefix(url, "http://") {
			url = url[7:]
		}

		data := UrlInfo{
			Message: message.Trailing,
			Nick:    message.Name,
			Time:    time.Now().Unix(),
		}
		b.handleUrl(bucket, url, &data)
	}
}

func (b *Bot) handleUrl(channel, url string, data *UrlInfo) {
	var oldData UrlInfo
	bucketName := BUCKET_PREFIX + channel

	err := b.db.Update(func(tx *bolt.Tx) (err error) {
		bucket := tx.Bucket([]byte(bucketName))

		// If bucket does not exist create one
		if bucket == nil {
			log.Println("No bucket found, creating:", bucketName)
			if bucket, err = tx.CreateBucket([]byte(bucketName)); err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			buf, _ := json.Marshal(data)
			err = bucket.Put([]byte(url), buf)
			return
		}

		v := bucket.Get([]byte(url))
		if v == nil {
			log.Println("No link found creating new:", url)
			buf, _ := json.Marshal(data)
			err = bucket.Put([]byte(url), buf)
			return
		}

		err = json.Unmarshal(v, &oldData)
		since := time.Now().Sub(time.Unix(oldData.Time, 0))
		b.irc.Privmsg(channel, fmt.Sprintf("\x02WANHA! %s linkkas tämän jo %s sitten", oldData.Nick, since))
		return
	})
	if err != nil {
		log.Println(err)
		return
	}

}
