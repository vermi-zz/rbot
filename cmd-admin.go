package main

import (
	irc "github.com/fluffle/goirc/client"
	"exec"
	"os"
	"time"
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
		cmd := exec.Command("./rbot")
		here, _ := os.Getwd()
		cmd.Dir = here
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			say(conn, channel, "Unable to start new process: %s", err)
			return
		}
		conn.Quit("Restarting.")
		time.Sleep(50)
		os.Exit(0)
	}
}
