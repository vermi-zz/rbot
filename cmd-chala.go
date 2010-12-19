package main

import (
	irc "github.com/fluffle/goirc/client"
	"http"
	"strings"
	"fmt"
	"io/ioutil"
	"json"
	"html"
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

func getq(conn *irc.Conn, nick *irc.Nick, id string, channel string) {
	if id == "" {
		say(conn, channel, "Syntax: !qdb ##; where ## is the unique ID of a quote in the database.")
		return
	}

	site := fmt.Sprintf("http://www.chalamius.se/quotes/api/json/quote/%s/", id)
	stuff, _, err := http.Get(site)
	defer stuff.Body.Close()

	if err != nil {
		say(conn, channel, "Something went wrong!")
		return
	}

	x, err := ioutil.ReadAll(stuff.Body)
	if err != nil {
		say(conn, channel, "Something went wrong!")
		return
	}

	var reply map[string]interface{}

	err = json.Unmarshal(x, &reply)
	if err != nil {
		say(conn, channel, "That quote ID doesn't seem to exist.")
		return
	}

	var text []string

	content := reply["content"].(string)
	content = html.UnescapeString(content)
	text = strings.Split(content, "\r\n", -1)
	who := reply["author"].(string)
	where := reply["channel"].(string)

	byline := "Submitted by " + who
	if where != "" {
		byline += " in " + where
	}

	if len(text) > 4 {
		say(conn, channel, "Quote %s: http://www.chalamius.se/quotes/quote.php?id=%s", id, id)
		return
	}

	for _, out := range text {
		say(conn, channel, "%s", out)
	}
	say(conn, channel, byline)
}
