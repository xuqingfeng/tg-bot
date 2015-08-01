package main

import (
	"fmt"
	"net/http"
	"os"
	"encoding/json"
	"strings"
	"heroku.com/tg-bot/util"
	"heroku.com/tg-bot/biubot"
	"heroku.com/tg-bot/google"
)

var hostname string
var config biubot.Config

func main() {

	config = biubot.GetConfig()

	http.HandleFunc("/", serveHomePage)

	webHookUrl := biubot.GetWebHookUrl()
	http.HandleFunc(webHookUrl, serveWebHook)

	http.HandleFunc("/google", serveGoogle)


	hostname, _ := os.Hostname()
	// auto change
	if hostname != "jsxqfdeMac-mini.local" {
		err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
		util.PanicIf(err)
	}else {
		err := http.ListenAndServe(":8080", nil)
		util.PanicIf(err)
	}

}

func serveWebHook(w http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	var update biubot.Update
	decoder.Decode(&update)

	// command
	var tgResponse string
	userInput := update.Message.Text
	commands := strings.Split(userInput, " ")
	if commands[0] == "/g" {
		query := strings.TrimLeft(userInput, "/g")
		links := google.SearchInGoogle(query)
		for i, link := range links {
			if ( i < config.Google_top) {
				tgResponse += link + "\n"
			}
		}
	} else if commands[0] == "/yo" {
		tgResponse = strings.TrimLeft(userInput, "/yo")
	} else if commands[0] == "/h" {
		tgResponse = "/g [query] - Search In Google\n"+
		"/yo - Take It Easy\n"+
		"/h  - List All Commands"
	} else {
		tgResponse = "Oops Master, Wrong Command!"
	}
	biubot.SendMessage(update.Message.From.Id, tgResponse)
}

func serveHomePage(w http.ResponseWriter, res *http.Request) {

	// hostname = res.Host
	// URL？
	hostname, _ := os.Hostname()
	fmt.Fprintln(w, hostname)
}



