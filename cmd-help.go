package main

import (
	irc "github.com/fluffle/goirc/client"
	config "goconfig"
	"strings"
)

const helpFile = "help.conf"

func helpGetTopics() (helpTopics string) {
	c, _ := config.ReadDefault(helpFile)
	helpTopics, _ = c.String("DEFAULT", "topics")
	return helpTopics
}

func helpGetHelp(conn *irc.Conn, topic string, channel string) {
	c, _ := config.ReadDefault(helpFile)

	var text []string
	content, _ := c.String(topic, "content")
	content = strings.Replace(content, "{trigger}", trigger, -1)
	text = strings.Split(content, "\n", -1)

	for _, out := range text {
		say(conn, channel, out)
	}
}

func helpProcessRequest(conn *irc.Conn, nick *irc.Nick, topic string, channel string) {
	channel = nick.Nick
	topic = strings.TrimSpace(topic)
	validTopics := strings.ToLower(helpGetTopics())

	if topic == "" {
		say(conn, channel, "Syntax is !help <topic>. Valid topics are:")
		say(conn, channel, validTopics)
	}

	topic = strings.ToLower(topic)

	if strings.Contains(validTopics, topic) {
		helpGetHelp(conn, topic, channel)
		return
	} else {
		say(conn, channel, "%s is not a valid help topic. Valid topics are:", topic)
		say(conn, channel, validTopics)
		return
	}
}
