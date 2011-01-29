package main

import (
	irc "github.com/fluffle/goirc/client"
	"strings"
	"time"
	config "goconfig"
	"strconv"
)

func op(conn *irc.Conn, nick *irc.Nick, args, target string) {
	channel, args := parseAccess(conn, nick, target, args, "o")
	if channel == "" {
		return
	}

	if args == "" {
		conn.Mode(channel, "+o " + nick.Nick)
	} else {
		ops := strings.TrimSpace(args)
		count := strings.Count(ops, " ") + 1
		modestring := "+" + strings.Repeat("o", count) + " " + ops
		conn.Mode(channel, modestring)
	}
}

func deop(conn *irc.Conn, nick *irc.Nick, args, target string) {
	channel, args := parseAccess(conn, nick, target, args, "o")
	if channel == "" {
		return
	}

	if args == "" {
		conn.Mode(channel, "-o " + nick.Nick)
	} else {
		ops := strings.TrimSpace(args)
		count := strings.Count(ops, " ") + 1
		modestring := "-" + strings.Repeat("o", count) + " " + ops
		conn.Mode(channel, modestring)
	}
}

func halfop(conn *irc.Conn, nick *irc.Nick, args, target string) {
	channel, args := parseAccess(conn, nick, target, args, "oh")
	if channel == "" {
		return
	}

	if args == "" {
		conn.Mode(channel, "+h " + nick.Nick)
	} else {
		// giving others +h requires o
		if !hasAccess(conn, nick, channel, "o") {
			return
		}
		halfops := strings.TrimSpace(args)
		count := strings.Count(halfops, " ") + 1
		modestring := "+" + strings.Repeat("h", count) + " " + halfops
		conn.Mode(channel, modestring)
	}
}

func dehalfop(conn *irc.Conn, nick *irc.Nick, args, target string) {
	channel, args := parseAccess(conn, nick, target, args, "oh")
	if channel == "" {
		return
	}

	if args == "" {
		conn.Mode(channel, "-h " + nick.Nick)
	} else {
		if !hasAccess(conn, nick, channel, "o") {
			return
		}
		halfops := strings.TrimSpace(args)
		count := strings.Count(halfops, " ") + 1
		modestring := "-" + strings.Repeat("h", count) + " " + halfops
		conn.Mode(channel, modestring)
	}
}

func voice(conn *irc.Conn, nick *irc.Nick, args, target string) {
	channel, args := parseAccess(conn, nick, target, args, "v")
	if channel == "" {
		return
	}

	if args == "" {
		conn.Mode(channel, "+v " + nick.Nick)
	} else {
		voices := strings.TrimSpace(args)
		count := strings.Count(voices, " ") + 1
		modestring := "+" + strings.Repeat("v", count) + " " + voices
		conn.Mode(channel, modestring)
	}
}

func devoice(conn *irc.Conn, nick *irc.Nick, args, target string) {
	channel, args := parseAccess(conn, nick, target, args, "v")
	if channel == "" {
		return
	}

	if args == "" {
		conn.Mode(channel, "-v " + nick.Nick)
	} else {
		voices := strings.TrimSpace(args)
		count := strings.Count(voices, " ") + 1
		modestring := "-" + strings.Repeat("v", count) + " " + voices
		conn.Mode(channel, modestring)
	}
}

func kick(conn *irc.Conn, nick *irc.Nick, args, target string) {
	channel, args := parseAccess(conn, nick, target, args, "oh")
	if channel == "" || args == "" {
		return
	}

	split := strings.Split(args, " ", 2)
	if n := conn.GetNick(split[0]); n == nil || (split[0] != nick.Nick &&
		(!hasAccess(conn, nick, channel, "o") && hasAccess(conn, n, channel, "oh"))) {
		// if we only have h, we can't kick people with o or h
		return
	}

	reason := "(" + nick.Nick + ")"
	if len(split) == 2 {
		reason += " " + split[1]
	}
	conn.Kick(channel, split[0], reason)
}

func ban(conn *irc.Conn, nick *irc.Nick, args, target string) {
	channel, args := parseAccess(conn, nick, target, args, "oh")
	if channel == "" || args == "" {
		return
	}

	bans := strings.TrimSpace(args)
	split := strings.Fields(bans)
	// turn nicks into *!*@host
	for i, ban := range(split) {
		if strings.Index(ban, "@") != -1 {
			// already a host
			continue
		}
		n := conn.GetNick(ban)
		if n == nil {
			//couldn't find the nick, so just cross our fingers
			continue
		}
		split[i] = "*!*@" + n.Host
	}
	bans = strings.Join(split, " ")
	modestring := "+" + strings.Repeat("b", len(bans)) + " " + bans
	conn.Mode(channel, modestring)
}

func unban(conn *irc.Conn, nick *irc.Nick, args, target string) {
	channel, args := parseAccess(conn, nick, target, args, "oh")
	if channel == "" || args == "" {
		return
	}

	ch := conn.GetChannel(channel)
	if ch == nil {
		say(conn, target , "%s: Unable to get channel information about %s", nick.Nick, channel)
		return
	}
	bans := strings.TrimSpace(args)
	split := strings.Fields(bans)
	for i, ban := range(split) {
		if strings.Index(ban, "@") != -1 {
			// it's already a host, do nothing
			continue
		}
		if b, ok := ch.Bans[ban]; ok {
			// we've seen this nick banned before
			split[i] = b
		} else if n := conn.GetNick(ban); n != nil {
			// the user is in one of our channels, here's our best guess
			split[i] = "*!*@" + n.Host
		} else if _, err := strconv.Atoi(ban); err == nil {
			// the ban is an integer, let's find it in the banlist
            c, err := config.ReadDefault("bans.list")
            if err != nil { return }

            host, err := c.String(channel, ban + ".host")
            if err == nil { split[i] = host }
			banLogDel(channel, ban)
        }
	}
	bans = strings.Join(split, " ")
	modestring := "-" + strings.Repeat("b", len(bans)) + " " + bans
	conn.Mode(channel, modestring)
}

func banLogDel(channel string, ban string) {
	c, err := config.ReadDefault("bans.list")
	if err != nil { return }

	c.AddOption(channel, ban + ".status", "REMOVED")
	c.WriteFile("bans.list", 0644, "Ban List")
}

func banLogAdd(host string, nick string, reason string, channel string) {
	c, _ := config.ReadDefault("bans.list")
	banCount, _ := c.Int(channel, "count")
	banCount += 1
	id := strconv.Itoa(banCount)
	banTime := time.LocalTime().Format("%m/%D/%y @ %H:%M")
	
	c.AddOption(channel, "count", id)
	c.AddOption(channel, id + ".nick", nick)
	c.AddOption(channel, id + ".host", host)
	c.AddOption(channel, id + ".reason", reason)
	c.AddOption(channel, id + ".time", banTime)
	c.AddOption(channel, id + ".status", "ACTIVE")
	c.WriteFile("bans.list", 0644, "Ban List")
}

func banList(conn *irc.Conn, nick *irc.Nick, args string, target string) {
	channel, args := parseAccess(conn, nick, target, args, "?")
	if channel == "" {
		say(conn, target, "Please specify a channel.")
		return
	}

	c, err := config.ReadDefault("bans.list")
	if err != nil {
		say(conn, channel, "No banlist file exists for %s. This is probably an error.", channel)
		return
	}

	banCount, _ := c.Int(channel, "count")
	
	if banCount == 0 {
		say(conn, channel, "There are no bans for %s.", channel)
		return
	}
	
	howMany, err := strconv.Atoi(args)
	if err != nil { howMany = 10 }

	if howMany == 0 {
		say(conn, channel, "Why would you ask me to show you nothing? Sheesh.")
		return
	 }
        if howMany < 0 || howMany > banCount { howMany = banCount }
	
	say(conn, nick.Nick, "There are a total of %v bans in the log for %s.", banCount, channel)
	say(conn, nick.Nick, "Displaying the last %v bans.", howMany)

	if banCount == howMany { howMany = 0 } else { howMany = banCount - howMany }
	
	for counter := banCount; counter > howMany; counter -= 1 {
		id := strconv.Itoa(counter)
		logNick, _ := c.String(channel, id + ".nick")
		logReason, _ := c.String(channel, id + ".reason")
		logTime, _ := c.String(channel, id + ".time")
		logStatus, _ := c.String(channel, id + ".status")
		
		say(conn, nick.Nick, "%v: %s [%s] | %s | %s", counter, logTime, logStatus, logNick, logReason)
	}
}

func kickban(conn *irc.Conn, nick *irc.Nick, args, target string) {
	channel, args := parseAccess(conn, nick, target, args, "oh")
	if channel == "" || args == "" {
		return
	}

	split := strings.Split(args, " ", 2)

	n := conn.GetNick(split[0])
	if n == nil || (split[0] != nick.Nick &&
		(!hasAccess(conn, nick, channel, "o") && hasAccess(conn, n, channel, "oh"))) {
		return
	}
	conn.Mode(channel, "+b *!*@" + n.Host)

	reason := "(" + nick.Nick + ")"
	if len(split) == 2 {
		reason += " " + split[1]
	}
	conn.Kick(channel, split[0], reason)
	banLogAdd("*!*@" + n.Host, split[0], reason, channel)
}

func topic(conn *irc.Conn, nick *irc.Nick, args, target string) {
	channel, args := parseAccess(conn, nick, target, args, "t")
	if channel == "" {
		return
	}
	section := conn.Network + " " + channel
	if args != "" {
		updateConf(section, "basetopic", args)
		conn.Topic(channel, args)
	} else {
		basetopic, _ := conf.String(section, "basetopic")
		say(conn, nick.Nick, "Basetopic: %s", basetopic)
	}
}
func appendtopic(conn *irc.Conn, nick *irc.Nick, args, target string) {
	channel, args := parseAccess(conn, nick, target, args, "t")
	if channel == "" {
		return
	}
	c := conn.GetChannel(channel)
	if c == nil {
		say(conn, target, "Error while getting channel information for %s", channel)
		return
	}

	section := conn.Network + " " + channel
	basetopic, _ := conf.String(section, "basetopic")
	if basetopic == "" || !strings.HasPrefix(c.Topic, basetopic) {
		basetopic = c.Topic
		say(conn, nick.Nick, "New basetopic: %s", basetopic)
		updateConf(section, "basetopic", basetopic)
	}
	conn.Topic(channel, basetopic + args)
}

func part(conn *irc.Conn, nick *irc.Nick, args, target string) {
	channel, args := parseChannel(target, args)
	if channel == "" {
		return
	}
	user := user(nick)
	if owner, _ := auth.String(conn.Network, "owner"); owner == user {
		conn.Part(channel, "")
	}
}
