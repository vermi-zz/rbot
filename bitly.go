package main

import (
	"http"
	"fmt"
	"json"
	"io/ioutil"
)

func shorten(long string) (short string) {
	key := "R_e659dbb5514e34edc3540a7c95b0041b"
	login := "jvermillion"

	long = http.URLEscape(long)

	url := fmt.Sprintf("http://api.bit.ly/v3/shorten?login=%s&apiKey=%s&longUrl=%s&format=json", login, key, long)

	r, _, err := http.Get(url)
        defer r.Body.Close()

        if err != nil {
                return "Error connecting to bit.ly"
        }

        b, err := ioutil.ReadAll(r.Body)
        if err != nil {
                return "Error reading bit.ly response"
        }

        var j map[string]interface{}

        err = json.Unmarshal(b, &j)
        if err != nil {
                return "Unable to shorten URL."
        }

	var data map[string]interface{} = j["data"].(map[string]interface{})

        return data["url"].(string)
}
