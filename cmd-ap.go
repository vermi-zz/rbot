package main

import(
        irc "github.com/fluffle/goirc/client"
	"github.com/kless/goconfig/config"
)

const configFile = "ap.cfg"

func readAPConfig(nick *irc.Nick, channel string) (apnick string) {
	c, _ := config.ReadDefault(configFile)
	
	hostmask := user(nick)

	apnick, _ = c.String(channel, hostmask)

	return apnick
}

func apProfile(conn *irc.Conn, nick *irc.Nick, arg string, channel string) {
	username := arg

	if username  == "" {
		username = readAPConfig(nick, channel)
	}

	if username == "" {
		username = nick.Nick
	}

	say(conn, channel, "%s's profile: http://anime-planet.com/users/%s", username, username)
}

func animelist(conn *irc.Conn, nick *irc.Nick, arg string, channel string) {
        username := arg

        if username == "" {
                username = readAPConfig(nick, channel)
        }

        if username == "" {
                username = nick.Nick
        }

	say(conn, channel, "%s's anime list: http://anime-planet.com/users/%s/anime", username, username)
}

func mangalist(conn *irc.Conn, nick *irc.Nick, arg string, channel string) {
        username := arg

        if username == "" {
                username = readAPConfig(nick, channel)
        }

        if username == "" {
                username = nick.Nick
        }

	say(conn, channel, "%s's manga list: http://anime-planet.com/users/%s/manga", username, username)
}

func apnick(conn *irc.Conn, nick *irc.Nick, arg string, channel string) {
	if arg == "" {
		say(conn, channel, "Format is !apnick <nickname>")
		return
	}

	c, _ := config.ReadDefault(configFile)

	hostmask := user(nick)

	c.AddOption(channel, hostmask, arg)
	c.WriteFile(configFile, 0644, "")

	say(conn, channel, "AP Nick set to %s", arg)
}
