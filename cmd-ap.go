package main

import (
	irc    "github.com/fluffle/goirc/client"
	config "goconfig"
	"http"
	"strings"
	"io/ioutil"
	"strconv"
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
	username := strings.TrimSpace(arg)

	if username == "" {
		username = apReadConfig(nick)
	}

	if username == "" {
		username = nick.Nick
	}

	if apUserExists(username) {
		say(conn, channel, "%s's profile: http://www.anime-planet.com/users/%s", username, username)
		return
	}

	say(conn, channel, "The user '%s' doesn't exist. Try again.", username)
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

func apStatsUID(nick string) int {
	url := "http://www.raylu.net/ap/user.php?nick=" + http.URLEscape(nick)
	r, _, err := http.Get(url)
	defer r.Body.Close()
	if err != nil || r.StatusCode != 200 {
		return -1
	}

	b, _ := ioutil.ReadAll(r.Body)
	uid, err := strconv.Atoi(string(b))
	if err != nil {
		return -1;
	}
	return uid
}
func apMyStats(conn *irc.Conn, nick *irc.Nick, arg string, channel string) {
	apStats(conn, nick, nick.Nick, channel)
}
func apStats(conn *irc.Conn, nick *irc.Nick, arg string, channel string) {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		say(conn, channel, "Channel stats: https://www.raylu.net/ap")
		return
	}

	uid := apStatsUID(arg)
	if uid >= 0 {
		say(conn, channel, "Stats for %s: https://www.raylu.net/ap/user.php?uid=%d", arg, uid)
	} else {
		say(conn, channel, "Could not find stats for %s", arg)
	}
}
