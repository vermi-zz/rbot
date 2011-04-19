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

func apStatsUID(nick string)(uid int){
		url := "https://www.raylu.net/ap/user.php?nick=" + http.URLEscape(nick)

	r, err := http.Head(url)

	if err != nil {
		return -1
	}

	if r.StatusCode == 200 {
		r, _, _  = http.Get(url)
		b, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		i, _ := strconv.Atoi(string(b))
		r.Body.Close()
		return i
	}

	r.Body.Close()
	return -1
}
func apMyStats(conn *irc.Conn, nick *irc.Nick, arg string, channel string){
	apStats(conn, nick, nick.Nick, channel)
}
func apStats(conn *irc.Conn, nick *irc.Nick, arg string, channel string){
	arg = strings.TrimSpace(arg)
	split := strings.Split(arg, " ", 2)
	arg = split[0]

	if arg == "" {
		say(conn, channel, "Channel Stats: https://www.raylu.net/ap")
		return
	}
		
		uid := apStatsUID(arg)
		if uid > 0 {
			say(conn, channel, "Stats for %s: https://www.raylu.net/ap/user.php?uid=%v", arg, uid)
		} else { say(conn, channel, "Channel stats: https://www.raylu.net/ap") }
}
