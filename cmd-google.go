package main

import (
	irc "github.com/fluffle/goirc/client"
	"fmt"
	"strings"
	"os"
	"http"
	"net"
	"url"
	"io/ioutil"
	"html"
	"regexp"
	"strconv"
	"utf8"
)

func tr(conn *irc.Conn, nick *irc.Nick, args, target string) {
	if args == "" {
		return
	}

	var sourcelang, targetlang, text string

	index := strings.IndexAny(args, " 　") // handle spaces and ideographic spaces (U+3000)
	if index == 5 && args[2] == '|' {
		sourcelang = args[:2]
		targetlang = args[3:5]
		if args[5] == ' ' {
			text = args[6:]
		} else {
			text = args[5 + utf8.RuneLen(3000):]
		}
	} else {
		sourcelang = "auto"
		targetlang = "en"
		text = args
	}

	say(conn, target, translate(sourcelang, targetlang, text))
}

func roman(conn *irc.Conn, nick *irc.Nick, args, target string) {
	if args == "" {
		return
	}

	var sourcelang, targetlang string
	if utf8.NewString(args).IsASCII() {
		sourcelang = "en"
	} else {
		sourcelang = "ja"
	}
	targetlang, _ = conf.String(conn.Network, "roman")
	if targetlang == "" {
		targetlang = "ja"
	}
	say(conn, target, translate(sourcelang, targetlang, args))
}

func translate(sourcelang, targetlang, text string) string {
	uri := fmt.Sprintf("http://translate.google.com/translate_a/t?client=t&hl=%s&sl=%s&tl=en-US&text=%s",
		targetlang, sourcelang, url.QueryEscape(text))

	b, err := getUserAgent(uri)
	if err != nil {
		return "Error while requesting translation"
	}

	result := strings.SplitN(string(b), `"`, 11)
	if len(result) < 11 {
		return "Error while parsing translation"
	}

	if sourcelang == "auto" {
		return result[9]
	}

	source := utf8.NewString(result[1])
	romanized := utf8.NewString(result[5])
	if romanized.RuneCount() > 0 {
		if sourcelang == "en" && !strings.Contains(text, " ") {
			// Google duplicates when there is only one source word
			source_left := source.Slice(0, source.RuneCount()/2)
			source_right := source.Slice(source.RuneCount()/2, source.RuneCount())
			romanized_left := romanized.Slice(0, romanized.RuneCount()/2)
			romanized_right :=romanized.Slice(romanized.RuneCount()/2, romanized.RuneCount())
			if (source_left == source_right &&
				strings.ToLower(romanized_left) == strings.ToLower(romanized_right)) {
				return fmt.Sprintf( "%s: %s", source_left, romanized_left)
			}
		}
		return fmt.Sprintf("%s: %s", source, romanized)
	}
	return source.String()
}

func calc(conn *irc.Conn, nick *irc.Nick, args, target string) {
	if args == "" {
		return
	}
	uri := fmt.Sprintf("http://www.google.com/ig/calculator?hl=en&q=%s", url.QueryEscape(args))

	b, err := getUserAgent(uri)
	if err != nil {
		say(conn, target, "%s: Error while requesting calculation", nick.Nick); return
	}

	re := regexp.MustCompile(`{lhs: "(.*)",rhs: "(.*)",error: "(.*)",icc: (true|false)}`)
	result := re.FindSubmatch(b)
	if len(result) != 5 {
		say(conn, target, "%s: Error while parsing.", nick.Nick)
		return
	}
	if len(result[3]) > 1 {
		output := fmt.Sprintf(`"%s"`, result[3])
		error := parseCalc(output)
		if error != "" {
			say(conn, target, "%s: Error: %s", nick.Nick, error)
		} else {
			say(conn, target, "%s: Error while calculating and error while decoding error.", nick.Nick)
		}
		return
	}
	if len(result[1]) == 0 || len(result[2]) == 0 {
		say(conn, target, "%s: Error while calculating.", nick.Nick)
		return
	}

	output := fmt.Sprintf(`"%s = %s"`, result[1], result[2])
	output = parseCalc(output)
	if output == "" {
		say(conn, target, "%s: Error while decoding.", nick.Nick); return
	}
	say(conn, target, output)
}

func parseCalc(output string) string {
	parsed, err := strconv.Unquote(output)
	if err != nil {
		return ""
	}
	parsed = html.UnescapeString(parsed)
	parsed = strings.Replace(parsed, "<sup>", "^(", -1)
	parsed = strings.Replace(parsed, "</sup>", ")", -1)
	return parsed
}

// make a GET request with a fake user agent
// this is definitely not for those undocumented Google APIs
func getUserAgent(urlstr string) ([]byte, os.Error) {
	urlobj, err := url.Parse(urlstr)
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial("tcp", urlobj.Host + ":http")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", urlstr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	httpconn := http.NewClientConn(conn, nil)
	response, err := httpconn.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	return b, err
}
