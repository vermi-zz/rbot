package main

import (
	"irc"
	"http"
	"strings"
)

func addquote(conn *irc.Conn, nick *irc.Nick, quote string, channel string) {
	if quote == "" {
		say(conn, channel, "Syntax: !quote <quote text>; use \\n as a newline separator.")
		say(conn, channel, "Example: !quote <vermi> this is a quote")
		return
	}

	data := map[string]string{
		"name": nick.Nick,
		"channel": channel,
		"content": strings.Replace(quote, "\\n", "\n", -1),
	}

	site := "http://www.chalamius.se/quotes/takesubmit.php"

	r, err := http.PostForm(site, data)
	r.Body.Close()

	if err != nil {
		say(conn, channel, "There was an error while posting.")
		return
	}

	say(conn, channel, "Your quote has been submitted.")
}
