package main

import (
	"irc"
	"fmt"
	"http"
	"xml"
)

func scrape(conn *irc.Conn, channel string, site string) (status string) {
	stuff, _, err := http.Get(site)
	if err != nil {
		return "DOWN"
	}

        type SecondResponse struct {
                File_url string "attr"
                Rating string "attr"
        }

        type FirstResponse struct {
                Response SecondResponse
        }

        type IbSearch struct {
                Response FirstResponse
        }

        var url IbSearch

        err = xml.Unmarshal(stuff.Body, &url)
        if err != nil {
       		return "FAIL"
        }

        rating := ""
        switch {
                case url.Response.Response.Rating == "s": rating = "[Rating: Safe] "
		case url.Response.Response.Rating == "q": rating = "[Rating: Questionable] "
		case url.Response.Response.Rating == "e": rating = "[Rating: Explicit] "
		case url.Response.Response.Rating == "0": rating = "[Not Rated] "
        }

        result := fmt.Sprintf("%s%s", rating, url.Response.Response.File_url)

        say(conn, channel, result)

	return "OK"
}

func booru(conn *irc.Conn, nick *irc.Nick, tag string, channel string) {
	if tag == "" {
		channel = nick.Nick
		say(conn, channel, "Syntax: !booru tag [tag2 tag3 ...]; you can use +tag, -tag, fav:, and rating:[s,q,e]")
		say(conn, channel, "Example: !booru +loli +pantsu -blonde* rating:s")
		say(conn, channel, "Results could take up to 20 seconds to appear.")
		return
	}

	tag = http.URLEscape(tag)
	site := fmt.Sprintf("http://www.i-forge.net/imageboards/?action=randimage&randimage[phrase]=%s&format=xml", tag)

	status := "FAIL"
	for fail := 0; fail < 10 && status == "FAIL"; fail++ {
		status = scrape(conn, channel, site)
	}

	switch {
		case status == "FAIL": say(conn, channel, "I looked and looked and just couldn't find anything. Try again in a bit.")
		case status == "DOWN": say(conn, channel, "booru search is down. If this keeps happening, please inform the owner.")
	}
}

func futa(conn *irc.Conn, nick *irc.Nick, tag string, channel string) {
	tag = "futa futanari -futaba*"
	tag = http.URLEscape(tag)
	site := fmt.Sprintf("http://www.i-forge.net/imageboards/?action=randimage&randimage[phrase]=%s&format=xml", tag)

	status := scrape(conn, channel, site)

        switch {
                case status == "FAIL": say(conn, channel, "%s, you are a pervert! No futa for you! >:|", nick.Nick)
                case status == "DOWN": say(conn, channel, "booru search is down. If this keeps happening, please inform the owner.")
        }
}

func loli(conn *irc.Conn, nick *irc.Nick, tag string, channel string) {
        tag = "loli*"
        tag = http.URLEscape(tag)
        site := fmt.Sprintf("http://www.i-forge.net/imageboards/?action=randimage&randimage[phrase]=%s&format=xml", tag)

        status := scrape(conn, channel, site)

        switch {
                case status == "FAIL": say(conn, channel, "%s, you are a pervert! No loli for you! >:|", nick.Nick)
                case status == "DOWN": say(conn, channel, "booru search is down. If this keeps happening, please inform the owner.")
        }
}

func sloli(conn *irc.Conn, nick *irc.Nick, tag string, channel string) {
        tag = "+loli* rating:s"
        tag = http.URLEscape(tag)
        site := fmt.Sprintf("http://www.i-forge.net/imageboards/?action=randimage&randimage[phrase]=%s&format=xml", tag)

        status := scrape(conn, channel, site)

        switch {
                case status == "FAIL": say(conn, channel, "Aww, I couldn't find anything this time... ;-;")
                case status == "DOWN": say(conn, channel, "booru search is down. If this keeps happening, please inform the owner.")
        }
}
