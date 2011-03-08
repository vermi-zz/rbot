package main

import (
	irc    "github.com/fluffle/goirc/client"
	config "goconfig"
	"http"
	"strings"
	"rand"
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

func roll(conn *irc.Conn, nick *irc.Nick, arg string, channel string) {
	x := 1
	y := 6

	split := strings.Split(arg, "d", 2)
	if len(split) != 2 {
		split = []string{"1", "6"}
	}

	x, err := strconv.Atoi(split[0])
	if err != nil {
		x = 1
	}
	y, err = strconv.Atoi(split[1])
	if err != nil {
		y = 6
	}

	results := []string{}
	total := 0

	for i := x; i > 0; i-- {
		random := rand.Intn(y - 1) + 1
		total += random
		results = append(results, strconv.Itoa(random))
	}

	tmp := strings.Join(results, ", ")

	if x > 10 {
		say(conn, channel, "%s rolls %dd%d for a total of %d", nick.Nick, x, y, total)
	} else {
		say(conn, channel, "%s rolls %dd%d: %s, Total: %d", nick.Nick, x, y, tmp, total)
	}
}
