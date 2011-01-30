package main

import (
	irc "github.com/fluffle/goirc/client"
	"exec"
	"os"
)

func nick(conn *irc.Conn, nick *irc.Nick, args, target string) {
	if len(args) == 0 {
		return
	}
	owner, _ := auth.String(conn.Network, "owner")
	if owner == user(nick) {
		conn.Nick(args)
	}
}

func csay(conn *irc.Conn, nick *irc.Nick, args, target string) {
	channel, args := parseAccess(conn, nick, target, args, "s")
	if len(channel) > 0 {
		say(conn, channel, args)
	}
}

func restart(conn *irc.Conn, nick *irc.Nick, args, channel string) {
	owner, _ := auth.String(conn.Network, "owner")
	if owner == user(nick) {
		here, _ := os.Getwd()
		_, err := exec.LookPath("./rbot")
		if err != nil {
			say(conn, channel, "Couldn't find binary.")
			return
		}
		argv := []string{""}
		envv := []string{""}
		_, err = exec.Run("rbot", argv, envv, here, exec.PassThrough, exec.PassThrough, exec.PassThrough)
		if err != nil {
			say(conn, channel, "Unable to start new process.")
			return
		}
		say(conn, channel, "Restarting.")
		os.Exit(0)
	}
}