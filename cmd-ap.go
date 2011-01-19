package main

import(
        irc "github.com/fluffle/goirc/client"
	"github.com/kless/goconfig/config"
	"http"
)

const configFile = "ap.conf"

func apUserExists(username string) (isuser bool) {
	url := "http://anime-planet.com/users/" + username

	r, _, err := http.Get(url)

	if err != nil {
		return false
	}

	if r.StatusCode == 200 {
		return true
	}

	r.Body.Close()

	return false
}

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

	if apUserExists(username) {
		say(conn, channel, "%s's profile: http://anime-planet.com/users/%s", username, username)
		return
	}
	
	say(conn, channel, "The user '%s' doesn't exist. Try setting your AP username with !apnick.", username)
}

func animelist(conn *irc.Conn, nick *irc.Nick, arg string, channel string) {
        username := arg

        if username == "" {
                username = readAPConfig(nick, channel)
        }

        if username == "" {
                username = nick.Nick
        }

	if apUserExists(username) {
		say(conn, channel, "%s's anime list: http://anime-planet.com/users/%s/anime", username, username)
		return
	}

	say(conn, channel, "The user '%s' doesn't exist. Try setting your AP username with !apnick.", username)
}

func mangalist(conn *irc.Conn, nick *irc.Nick, arg string, channel string) {
        username := arg

        if username == "" {
                username = readAPConfig(nick, channel)
        }

        if username == "" {
                username = nick.Nick
        }

	if apUserExists(username) {
		say(conn, channel, "%s's manga list: http://anime-planet.com/users/%s/manga", username, username)
		return
	}

	say(conn, channel, "The user '%s' doesn't exist. Try setting your AP username with !apnick.", username)
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
