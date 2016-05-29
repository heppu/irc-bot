# IRC-bot
This is IRC bot written in golang. It has two modules; ops and wanha.

## ops
This module keeps track how many ops each user has (incremental credentials).

This:
```
/MODE +o heppu
/MODE +o heppu
```
Would show that user heppu has ops count of 2.
### Command
- !ops - show you own ops count
- !ops nick1 nick2 - show ops count given for users.

## wanha
This module saves all links postetd to channels and notifies you with annoying message if you posted a link that someone else has already posted.

## Libaries used
[bolt](https://github.com/boltdb/bolt)
[irc](https://github.com/sorcix/irc)
[xurls](https://github.com/mvdan/xurls)
[irc-client](https://github.com/heppu/jun) Forked from [FSX](https://github.com/FSX/jun)
