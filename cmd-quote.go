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

func quoteAdd(conn *irc.Conn, nick *irc.Nick, quote string, channel string) {
	if quote == "" {
		say(conn, channel, "Syntax: !quote <quote text>; use \\n as a newline separator.")
		say(conn, channel, "Example: !quote <vermi> this is a quote")
		return
	}

	data := map[string]string{
		"name":    nick.Nick,
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

func quoteGet(conn *irc.Conn, nick *irc.Nick, id string, channel string) {
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

func quoteRand(conn *irc.Conn, nick *irc.Nick, _, channel string) {
	url := "http://www.chalamius.se/quotes/api/json/random/"
	r, _, err := http.Get(url)
	defer r.Body.Close()

	if err != nil {
		say(conn, channel, "Error connecting to QDB.")
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		say(conn, channel, "Error reading JSON response.")
		return
	}

	var reply map[string]interface{}

	err = json.Unmarshal(b, &reply)
	if err != nil {
		say(conn, channel, "Error fetching random quote.")
		return
	}

	var text []string

	content := reply["content"].(string)
	content = html.UnescapeString(content)
	text = strings.Split(content, "\r\n", -1)
	id := reply["id"].(string)
	who := reply["author"].(string)
	where := reply["channel"].(string)

	byline := "Submitted by " + who
	if where != "" {
		byline += " in " + where
	}

	if len(text) > 4 {
		say(conn, channel, "Random quote: http://www.chalamius.se/quotes/quote.php?id=%s", id)
		return
	}

	for _, out := range text {
		say(conn, channel, "%s", out)
	}

	say(conn, channel, byline)
}

func quoteSearch(conn *irc.Conn, nick *irc.Nick, term string, channel string) {
	channel = nick.Nick
	term = strings.TrimSpace(term)

	if term == "" {
		say(conn, channel, "No search term specified. Defaulting to random quote.")
		quoteRand(conn, nick, term, channel)
		return
	}

	url := fmt.Sprintf("http://www.chalamius.se/quotes/api/json/search/%s", term)

	r, _, err := http.Get(url)
	defer r.Body.Close()

	if err != nil {
		say(conn, channel, "Error connecting to QDB.")
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		say(conn, channel, "Error reading JSON response.")
		return
	}

	var reply []map[string]interface{}

	err = json.Unmarshal(b, &reply)
	if err != nil {
		say(conn, channel, "No results. Try a different search term.")
		return
	}

	var searchResult map[string]interface{}
	var resultUrl string
	var i int

S:
	for i, searchResult = range reply {
		count := i + 1
		resultUrl = fmt.Sprintf("http://www.chalamius.se/quotes/quote/%s/", searchResult["id"].(string))
		say(conn, channel, "Result %v: %s (or !gq %s)", count, resultUrl, searchResult["id"].(string))

		switch i {
		case 4:
			break S
		}
	}

	if i == 4 {
		say(conn, channel, "Top 5 results returned.")
	} else {
		say(conn, channel, "All results returned.")
	}
}
