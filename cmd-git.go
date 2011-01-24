package main

import (
	"http"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"strings"
	"json"
	"io/ioutil"
	"strconv"
)
var gitProject = readConfString("DEFAULT", "git_project")
var gitLogin = readConfString("DEFAULT", "git_login")
var gitToken = readConfString("DEFAULT", "git_token")

func gitOpenIssue(conn *irc.Conn, nick *irc.Nick, body string, channel string) {
	if body == "" {
		say(conn, channel, "Syntax: !bug <body>; newlines should be separated with \\n.")
		return
	}

	if len(body) < 30 {
		say(conn, channel, "Error: to prevent spam, bug reports must be longer than 30 characters.")
		return
	}

	data := map[string]string{
		"login": gitLogin,
		"token": gitToken,
		"title": "IRC Issue from " + nick.Nick + " (" + channel + ")",
		"body": strings.Replace(body, "\\n", "\n", -1),
	}

	url := fmt.Sprintf("http://github.com/api/v2/json/issues/open/%s/%s", gitLogin, gitProject)

	r, err := http.PostForm(url, data)
	defer r.Body.Close()

	if err != nil {
		say(conn, channel, "Could not submit issue.")
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		say(conn, channel, "Error getting JSON response.")
		return
	}

	var reply map[string]interface{}

	err = json.Unmarshal(b, &reply)
	if err != nil {
		say(conn, channel, "Error parsing JSON response.")
		return
	}

	var issue map[string]interface{} = reply["issue"].(map[string]interface{})

	id := strconv.Ftoa64(issue["number"].(float64), 'f', -1)
	
	say(conn, channel, "Your issue has been submitted. You can view it at http://github.com/%s/rbot/issues/#issue/%s", gitLogin, id)
}
