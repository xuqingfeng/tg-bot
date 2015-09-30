package main

import (
	"encoding/json"
	"heroku.com/tg-bot/biubot"
	"heroku.com/tg-bot/common"
	"heroku.com/tg-bot/firebase"
	"heroku.com/tg-bot/google"
	"heroku.com/tg-bot/qiniu"
	"heroku.com/tg-bot/reddit"
	"heroku.com/tg-bot/util"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var config common.Config

func main() {

	config = common.GetConfig()

	http.HandleFunc("/", serveHomePage)
	http.Handle("/assets/", http.FileServer(http.Dir("templates/")))

	webHookUrl := biubot.GetWebHookUrl()
	http.HandleFunc(webHookUrl, serveWebHook)

	hostname, _ := os.Hostname()
	// auto change
	if hostname != "jsxqfdeMac-mini.local" {
		err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
		util.PanicIf(err)
	} else {
		err := http.ListenAndServe(":8080", nil)
		util.PanicIf(err)
	}

}

func serveWebHook(w http.ResponseWriter, req *http.Request) {

	firebase.IncRequest()
	decoder := json.NewDecoder(req.Body)
	var update biubot.Update
	decoder.Decode(&update)

	// command
	var tgResponse string
	userInput := update.Message.Text
	location := update.Message.Location
	photo := update.Message.Photo

	if userInput != "" {
		commands := strings.Split(userInput, " ")
		if commands[0] == "/g" {
			query := strings.TrimLeft(userInput, "/g")
			links := google.SearchInGoogle(query)
			for i, link := range links {
				if i < config.GoogleTop {
					tgResponse += link + "\n"
				}
			}
		} else if commands[0] == "/yo" {

			randNum := rand.Intn(10)
			photoLinks := reddit.GetPhoto()
			tgResponse = photoLinks[randNum]
		}else if commands[0] == "/i" {

			tgResponse = "https://github.com/xuqingfeng/tg-bot"
		}else if commands[0] == "/h" {
			tgResponse = "/g [query] - Search In Google\n" +
			"/yo - Send Me A Picture\n" +
			"/i - About Me\n" +
			"/h - List All Commands\n" +
			"[photo] Send It To Qiniu & Return The Link"
		} else {
			tgResponse = "Oops Master, Wrong Command!"
		}

		biubot.SendMessage(update.Message.From.Id, tgResponse)

	} else if location != (biubot.Location{}) {

		tgResponse = "Location: longitude:" + strconv.FormatFloat(location.Longitude, 'f', 6, 64) + " latitude:" + strconv.FormatFloat(location.Latitude, 'f', 6, 64)
		biubot.SendMessage(update.Message.From.Id, tgResponse)
	} else if photo != nil {

		maxKey := 0
		for k, photoSize := range photo {
			tgResponse += photoSize.File_id + " " + strconv.Itoa(photoSize.Width) + " " + strconv.Itoa(photoSize.Height) + " " + strconv.Itoa(photoSize.File_size) + "\n"
			maxKey = k
		}

		filePath := biubot.GetFile(photo[maxKey].File_id)
		if filePath != "" {
			r := biubot.GetFileData(filePath)
			ok := qiniu.UploadFile(photo[maxKey].File_id, photo[maxKey].File_size, r)
			if ok {
				tgResponse = config.QiniuDomain + photo[maxKey].File_id
			}else {
				tgResponse = "Error"
			}
		} else {
			tgResponse = "Error"
		}

		biubot.SendMessage(update.Message.From.Id, tgResponse)
	}

}

func serveHomePage(w http.ResponseWriter, res *http.Request) {

	RequestNum := firebase.GetRequest()
	tp, _ := template.ParseFiles("templates/index.html")
	tp.Execute(w, RequestNum)
}
