package main

import (
	irc    "github.com/fluffle/goirc/client"
	config "goconfig"
	"http"
	"strings"
)

const apConfigFile = "ap.conf"

func apUserExists(username string) (isuser bool) {
	url := "http://www.anime-planet.com/users/" + http.URLEscape(username)

	r, err := http.Head(url)

	if err != nil {
		return false
	}

	if r.StatusCode == 200 {
		return true
	}

	r.Body.Close()

	return false
}

func apReadConfig(nick *irc.Nick) (apnick string) {
	c, _ := config.ReadDefault(apConfigFile)

	hostmask := user(nick)

	apnick, _ = c.String(hostmask, "nick")

	return apnick
}

func apProfile(conn *irc.Conn, nick *irc.Nick, arg string, channel string) {
	var query_nick *irc.Nick
	if arg == "" {
		query_nick = nick
	} else {
		arg = strings.TrimSpace(arg)
		query_nick = conn.GetNick(arg)
		if query_nick == nil {
			say(conn, channel, "Could not find nick %s", arg)
			return
		}
	}

	username := apReadConfig(query_nick)
	if username == "" {
		username = query_nick.Nick
		if apUserExists(username) {
			say(conn, channel, "%s's profile: http://www.anime-planet.com/users/%s", username, username)
		} else {
			say(conn, channel, "The user '%s' doesn't exist. Try again.", username)
		}
	} else {
		say(conn, channel, "%s's profile: http://www.anime-planet.com/users/%s", query_nick.Nick, username)
	}
}

func apAnimeList(conn *irc.Conn, nick *irc.Nick, arg string, channel string) {
	username := strings.TrimSpace(arg)

	if username == "" {
		username = apReadConfig(nick)
	}

	if username == "" {
		username = nick.Nick
	}

	if apUserExists(username) {
		say(conn, channel, "%s's anime list: http://www.anime-planet.com/users/%s/anime", username, username)
		return
	}

	say(conn, channel, "The user '%s' doesn't exist. Try again.", username)
}

func apMangaList(conn *irc.Conn, nick *irc.Nick, arg string, channel string) {
	username := strings.TrimSpace(arg)

	if username == "" {
		username = apReadConfig(nick)
	}

	if username == "" {
		username = nick.Nick
	}

	if apUserExists(username) {
		say(conn, channel, "%s's manga list: http://www.anime-planet.com/users/%s/manga", username, username)
		return
	}

	say(conn, channel, "The user '%s' doesn't exist. Try again.", username)
}

func apSetNick(conn *irc.Conn, nick *irc.Nick, arg string, channel string) {
	arg = strings.TrimSpace(arg)

	if arg == "" {
		say(conn, channel, "Format is !apnick <nickname>")
		return
	}

	if apUserExists(arg) {
		c, _ := config.ReadDefault(apConfigFile)

		hostmask := user(nick)

		c.AddOption(hostmask, "nick", arg)
		c.WriteFile(apConfigFile, 0644, "")

		say(conn, channel, "Your anime-planet.com username has been recorded as '%s'", arg)
		return
	}

	say(conn, channel, "The user '%s' doesn't exist. Try again.", arg)
}

func apMyNick(conn *irc.Conn, nick *irc.Nick, _, channel string) {
	username := apReadConfig(nick)

	if username == "" {
		say(conn, channel, "You haven't set your username. You can do so with !apnick <username>")
		return
	}

	say(conn, channel, "Your anime-planet.com username has been recorded as '%s'.", username)
}

func apMyStats(conn *irc.Conn, nick *irc.Nick, arg string, channel string) {
	apStats(conn, nick, nick.Nick, channel)
}
func apStats(conn *irc.Conn, nick *irc.Nick, arg string, channel string) {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		say(conn, channel, "Channel stats: http://www.raylu.net/irc/ap.html")
	} else {
		say(conn, channel, "Stats for %s: http://www.raylu.net/irc/user.php?cid=ap&nick=%s", arg, arg)
	}
}
