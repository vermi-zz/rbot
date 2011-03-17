package main

import (
	irc "github.com/fluffle/goirc/client"
	"strings"
	"rand"
	"strconv"
)

func roll(conn *irc.Conn, nick *irc.Nick, arg string, channel string) {
	x := 1
	y := 6

	arg = strings.ToLower(arg)

	split := strings.Split(arg, "d", 2)
	if len(split) != 2 {
		split = []string{"1", "6"}
	}

	x, err := strconv.Atoi(split[0])
	if err != nil {
		x = 1
	}
	if x > 100 {
		x = 100
	}
	if x <= 0 {
		x = 1
	}
	y, err = strconv.Atoi(split[1])
	if err != nil {
		y = 6
	}
	if y <= 0 {
		y = 6
	}

	results := []string{}
	total := 0

	for i := x; i > 0; i-- {
		random := rand.Intn(y-1) + 1
		total += random
		results = append(results, strconv.Itoa(random))
	}

	tmp := strings.Join(results, ", ")

	if x > 10 {
		say(conn, channel, "%s rolls %dd%d for a total of %d", nick.Nick, x, y, total)
	} else {
		say(conn, channel, "%s rolls %dd%d: %s, Total: %d", nick.Nick, x, y, tmp, total)
	}
}

func eightBall(conn *irc.Conn, nick *irc.Nick, arg string, channel string) {
	if arg == "" {
		say(conn, channel, "You didn't ask anything!")
	} else {
		r := rand.Intn(19)
		s := ""
		switch r {
		case 0:
			s = "As I see it, yes."
		case 1:
			s = "It is certain."
		case 2:
			s = "It is decidedly so."
		case 3:
			s = "Most likely."
		case 4:
			s = "Outlook good."
		case 5:
			s = "Signs point to yes."
		case 6:
			s = "Without a doubt."
		case 7:
			s = "Yes."
		case 8:
			s = "Yes - definitely."
		case 9:
			s = "You may rely on it."
		case 10:
			s = "Reply hazy, try again."
		case 11:
			s = "Ask again later."
		case 12:
			s = "Better not tell you now..."
		case 13:
			s = "Cannot predict now."
		case 14:
			s = "Concentrate, then ask again."
		case 15:
			s = "Don't count on it."
		case 16:
			s = "My reply is no."
		case 17:
			s = "My sources say no."
		case 18:
			s = "Outlook not so good."
		case 19:
			s = "Very doubtful."
		}

		say(conn, channel, "%s: %s", nick.Nick, s)
	}
}
