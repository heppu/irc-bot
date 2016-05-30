# IRC-bot
IRC bot written in Golang which has two modules:

## ops
This module keeps track of how many ops each user has (incremental credentials).

This:
```
/MODE +o heppu
/MODE +o heppu
```
would show that user heppu has ops count of 2.
##### Commands
- !ops - show your ops count
- !ops nick1 nick2 - show ops count for given users

## wanha
This module saves all links posted to channels. It notifies you with an annoying message if you posted a link that someone else has already posted.

## Libraries used
- [boltdb](https://github.com/boltdb/bolt) Key value store written in Go
- [sorcix irc](https://github.com/sorcix/irc) Simple irc library
- [xurls](https://github.com/mvdan/xurls) Url parser
- [jun irc-client](https://github.com/heppu/jun) IRC client, forked from [FSX](https://github.com/FSX/jun)
