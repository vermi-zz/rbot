package main

import (
	irc "github.com/fluffle/goirc/client"
	"fmt"
	"os"
	"strings"
	"time"
	"crypto/tls"
	"crypto/rand"
	"goconfig"
	"container/list"
	"strconv"
)

const confFile = "rbot.conf"
var trigger string
var sections []string
var conf *config.Config

func main() {
	readConf()
	trigger = readConfString("DEFAULT", "trigger")
	readAuth()

	sections = conf.Sections()
	for _, s := range sections {
		if strings.Index(s, " ") == -1 && s != "DEFAULT" {
			// found a network
			go connect(s)
		}
	}

	<- make(chan bool)
}

func connect(network string) {
	if !readConfBool(network, "autoconnect") {
		return
	}
	server := readConfString(network, "server")
	nick := readConfString(network, "nick")
	user := readConfString(network, "user")
	nickserv, _ := conf.String(network, "nickserv")

	c := irc.New(nick, user, user)
	c.Network = network
	c.SSL = readConfBool(network, "ssl")
	if c.SSL {
		// we don't care about certificate validity
		c.SSLConfig = &tls.Config{Rand: rand.Reader, Time: time.Nanoseconds}
	}

	c.AddHandler("connected",
		func(conn *irc.Conn, line *irc.Line) {
			fmt.Printf("Connected to %s!\n", conn.Host)

			if len(nickserv) > 0 {
				conn.Privmsg("NickServ", "IDENTIFY " + nickserv)
			} else {
				autojoin(conn)
			}
		})
	c.AddHandler("privmsg", handlePrivmsg)
	c.AddHandler("mode", handleMode)
	c.AddHandler("join", handleJoin)
	c.AddHandler("invite", handleInvite)

	for {
		fmt.Printf("Connecting to %s...\n", server)
		if err := c.Connect(server); err != nil {
			fmt.Printf("Connection error: %s\n", err)
			break
		}
		for err := range c.Err {
			fmt.Printf("goirc error: %s\n", err)
		}
	}
}

func autojoin(conn *irc.Conn) {
	for _, s := range sections {
		split := strings.Split(s, " ", 2)
		if len(split) == 2 && split[0] == conn.Network {
			// found a channel
			if readConfBool(s, "autojoin") {
				fmt.Printf("Joining %s on %s\n", split[1], conn.Network)
				conn.Join(split[1])
			}
		}
	}
	go BanManager(conn)
}

func readConf() {
	var err os.Error
	conf, err = config.ReadDefault("rbot.conf")
	if (err != nil) {
		fmt.Printf("Config error: %s\n", err)
		os.Exit(1)
	}
}
func readConfString(section, option string) string {
	value, err := conf.String(section, option)
	if err != nil {
		panic(fmt.Sprintf("Config error: %s", err))
	}
	return value
}
func readConfBool(section, option string) bool {
	value, err := conf.Bool(section, option)
	if err != nil {
		panic(fmt.Sprintf("Config error: %s", err))
	}
	return value
}
func updateConf(section, option, value string) {
	conf.AddOption(section, option, value)
	if err := conf.WriteFile(confFile, 0644, ""); err != nil {
		panic("Error while writing to " + confFile)
	}
	// config.WriteFile destroys the config, so
	readConf()
}
func BanManager(conn *irc.Conn) {
	c, err := config.ReadDefault("bans.list")
	if err != nil {
		panic(fmt.Sprintf("Config error: %s", err))
	}
	if !c.HasOption("timed", "count") {
		BanManager(conn)
		return
	}
	count, err := c.Int("timed", "count")
	if err != nil || count == 0 {
		BanManager(conn)
		return
	}
	banlist := list.New()
	for i := count; i > 0; i++ {
		ban, _ := c.String("timed", strconv.Itoa(count))
		_ = banlist.PushBack(ban)
	}
	for e := banlist.Front(); e != nil; e.Next() {
		split := strings.Fields(e.Value.(string))
		expiry, _ := strconv.Atoi64(split[2])
		if expiry <= time.Seconds() {
			host, _ := c.String(split[0], split[1] + ".host")
			conn.Mode(split[0], "-b " + host)
			banLogDel(split[0], split[1])
			c.RemoveOption("timed", strconv.Itoa(count))
			count--
			c.AddOption("timed", "count", strconv.Itoa(count))
			c.AddOption(split[0], split[1] + ".status", "EXPIRED")
		}
	}
	c.WriteFile("bans.list", 0644, "Ban List")
	time.Sleep(300000000000)
	BanManager(conn)
}