package main

import (
	irc "github.com/fluffle/goirc/client"
	"fmt"
	"http"
	"url"
	"xml"
)

func booruDoSearch(conn *irc.Conn, channel string, site string) (status string) {
	stuff, err := http.Get(site)
	if err != nil {
		return "DOWN"
	}

	type SecondResponse struct {
		File_url string `xml:"attr"`
		Rating   string `xml:"attr"`
	}

	type FirstResponse struct {
		Response SecondResponse
	}

	type IbSearch struct {
		Response FirstResponse
	}

	var url_ IbSearch

	err = xml.Unmarshal(stuff.Body, &url_)
	if err != nil {
		return "FAIL"
	}

	rating := ""
	switch {
	case url_.Response.Response.Rating == "s":
		rating = "[Supposedly Safe] "
	case url_.Response.Response.Rating == "q":
		rating = "[Questionable] "
	case url_.Response.Response.Rating == "e":
		rating = "[Explicit] "
	case url_.Response.Response.Rating == "0":
		rating = "[Not Rated] "
	}

	result := rating + shorten(url_.Response.Response.File_url)

	say(conn, channel, result)

	return "OK"
}

func booruSearch(conn *irc.Conn, nick *irc.Nick, tag string, channel string) {
	if tag == "" {
		channel = nick.Nick
		say(conn, channel, "Syntax: !booru tag [tag2 tag3 ...]; you can use +tag, -tag, and rating:[s,q,e]")
		say(conn, channel, "Example: !booru +loli +pantsu -blonde* rating:s")
		say(conn, channel, "Results could take up to 20 seconds to appear.")
		return
	}

	tag = url.QueryEscape(tag)
	site := fmt.Sprintf("http://ibsearch.i-forge.net/?action=randimage&randimage[phrase]=%s&format=xml", tag)

	status := "FAIL"
	for fail := 0; fail < 10 && status == "FAIL"; fail++ {
		status = booruDoSearch(conn, channel, site)
	}

	switch {
	case status == "FAIL":
		say(conn, channel, "I looked and looked and just couldn't find anything. Try again in a bit.")
	case status == "DOWN":
		say(conn, channel, "booru search is down. If this keeps happening, please inform the owner.")
	}
}

func booruFuta(conn *irc.Conn, nick *irc.Nick, tag string, channel string) {
	tag = "futa futanari -futaba*"
	tag = url.QueryEscape(tag)
	site := fmt.Sprintf("http://ibsearch.i-forge.net/?action=randimage&randimage[phrase]=%s&format=xml", tag)

	status := booruDoSearch(conn, channel, site)

	switch {
	case status == "FAIL":
		say(conn, channel, "%s, you are a pervert! No futa for you! >:|", nick.Nick)
	case status == "DOWN":
		say(conn, channel, "booru search is down. If this keeps happening, please inform the owner.")
	}
}

func booruLoli(conn *irc.Conn, nick *irc.Nick, tag string, channel string) {
	tag = "loli*"
	tag = url.QueryEscape(tag)
	site := fmt.Sprintf("http://ibsearch.i-forge.net/?action=randimage&randimage[phrase]=%s&format=xml", tag)

	status := booruDoSearch(conn, channel, site)

	switch {
	case status == "FAIL":
		say(conn, channel, "%s, you are a pervert! No loli for you! >:|", nick.Nick)
	case status == "DOWN":
		say(conn, channel, "booru search is down. If this keeps happening, please inform the owner.")
	}
}

func booruSafeLoli(conn *irc.Conn, nick *irc.Nick, tag string, channel string) {
	tag = "+loli* rating:s"
	tag = url.QueryEscape(tag)
	site := fmt.Sprintf("http://ibsearch.i-forge.net/?action=randimage&randimage[phrase]=%s&format=xml", tag)

	status := booruDoSearch(conn, channel, site)

	switch {
	case status == "FAIL":
		say(conn, channel, "Aww, I couldn't find anything this time... ;-;")
	case status == "DOWN":
		say(conn, channel, "booru search is down. If this keeps happening, please inform the owner.")
	}
}
