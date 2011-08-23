rbot
======================

### Getting started

Assuming you have go set up (http://golang.org/),

	git clone git://github.com/vermi/rbot.git
	cd rbot
	git submodule init
	git submodule update
	gomake

rbot.conf and auth.conf will be copied. Configure those. bans.list will also be touched; it requires no configuration, but the bot won't run without it.
Afterwards, run the bot with the following command:

	./rbot

If you're going to be running the bot for extended periods of time, we recommended using the `nohup` command.
	nohup ./rbot &
This will keep the bot from stopping if you log off.

As a side note, cmd-ap.go is specific to #anime-planet.com on Rizon; it's not necessary for any normal functions of the bot, so you can remove it if you want. Just be sure to remove it from Makefile and handler.go as well.

### Flags

Access is configured in auth.conf, is per-channel, based on ident and host; nick is ignored. The owner is configured per server and other access is configured per channel. Owners can use any commands. Ignores are checked before access.

The following is a description of the commands enabled by each flag:

- `a`: add remove ignore unignore
- `o`: op halfop deop dehalfop kick|k ban|b unban|u kb
- `h`: halfop|hop dehalfop|dehop kick|k ban|b unban|u kb (hop and dehop can only be used on yourself and you cannot kick or kb people with o or h)
- `v`: voice, devoice
- `t`: topic appendtopic
- `s`: say

In addition, a user must have at least one flag to use `flags` and `list` (so users without access can't spam the bot).

### Commands

All commands are prefixed with the trigger configured in rbot.conf.

Access related commands:

- `flags raylu`: get's raylu's flags
- `flags`: get's your flags
- `add raylu t`: gives raylu the t flag
- `remove raylu t`: removes the t flag from raylu
- `remove raylu`: removes all of raylu's flags
- `ignore raylu`: ignores all messages from raylu's host on this network
- `unignore raylu`: unignore raylu
- `list`: lists all users and their access
- `list raylu`: searches the access list for raylu

Admin commands:

- `nick john`: changes the bot's nick to john
- `say text`: says text to the channel
- `bans`: Lists the ban-log for the current channel, if it exists.

The ban-log may need to be zeroed out manually from time to time to reduce clutter.

Op commands:

- `halfop|hop`: halfop yourself
- `halfop|hop raylu john`: halfop raylu and john
- `op`: op yourself
- `op raylu john`: op raylu and john
- `deop`: deop yourself
- `deop raylu john`: deop raylu and john
- `dehalfop`: dehalfop yourself
- `dehalfop|dehop raylu john`: dehalfop raylu and john
- `voice`: voice yourself
- `voice raylu john`: voice raylu and john
- `devoice`: devoice yourself
- `devoice raylu john`: devoice raylu and john
- `kick|k raylu`: kick raylu
- `ban|b raylu john!*@*`: ban raylu by hostname and john by nick
- `unban|u raylu john!*@*`: unban raylu by hostname and john by nick
- `kb raylu`: kick raylu, then ban him by hostname
- `topic text`: sets the topic and basetopic to text
- `topic`: gets the current basetopic
- `appendtopic text`: if the topic does not starts with basetopic, sets the basetopic to the current topic. Makes the topic basetopic+text.
- `part`: leaves the channel

Google API commands:

- `tr text`: detect the language of text
- `tr en|ja en|es text`: translate text into Japanese and Spanish
- `roman text`: translate text into romaji (see rbot.conf.example)
- `calc 1 usd in yen`: convert 1 USD to Japanese yen

Help commands:

- `help`: displays a list of possible help topics
- `help topic`: displays help on topic, as long as it's in the list of possible topics.

Help topics and content are manually entered in help.conf. You probably won't need to change this, unless you add commands to the bot yourself.

Imageboard Search Commands:

- `booru tag [tag2 tag3 ...]`: search the boorus for images containing any of the listed tags.
	* `+tag`: marks a tag as mandatory
	* `-tag`: marks a tag as undesired
	* `rating:[s,q,e]`: looks for images rated `s`afe, `q`uestionable, or `e`xplicit
- `loli`: returns a random loli image
- `sloli`: returns a random work-safe loli image
- `futa`: returns a random futa image

Anime-planet.com Commands:

- `apnick vermi`: records my site username as vermi (useful if I change nicks a lot)
- `mynick`: checks if I've set my site username, and tells me what it's set as
- `profile`: if I've set apnick, uses that value to return a link to my profile; otherwise, uses my current nick
- `profile vermi`: displays a link to vermi's profile; checks to see if a user named 'vermi' exists, first.
- `mymanga` and `myanime`: work the same as `profile`, but return links to Manga and Anime lists, respectively

Commands that don't require access behave the same when sent to a channel the bot is in and when whispered to the bot.

Commands that require access are listed above as if they were sent to a channel. When sent as a whisper, the first argument must be a channel name.

The bot will accept invites from the owner to any channel and only the owner can use `part`.

### Miscellaneous

This project was forked from jessta/goirc which is in turn forked from fluffle/goirc. Both of those projects are focused on developing the goirc framework whereas this is focused on developing a bot.
