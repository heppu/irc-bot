# IRC-bot
IRC bot written in golang which has two module:

## ops
This module keeps track how many ops each user has (incremental credentials).

This:
```
/MODE +o heppu
/MODE +o heppu
```
Would show that user heppu has ops count of 2.
##### Commands
- !ops - show your ops count
- !ops nick1 nick2 - show ops count for given users

## wanha
This module saves all links postetd to channels and notifies you with annoying message if you posted a link that someone else has already posted.

## Libaries used
- [boltdb](https://github.com/boltdb/bolt) Key value strore written in Go
- [sorcix irc](https://github.com/sorcix/irc) Simple irc library
- [xurls](https://github.com/mvdan/xurls) Url parser
- [jun irc-client](https://github.com/heppu/jun) IRC client, forked from [FSX](https://github.com/FSX/jun)
